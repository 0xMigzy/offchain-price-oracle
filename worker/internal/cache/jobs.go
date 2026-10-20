package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const jobTTL = 7 * 24 * time.Hour
const recentJobsKey = "oracle:jobs:recent"
const QueueStream = "oracle:jobs:stream"

var ErrJobNotFound = errors.New("job not found")

type Event struct {
	Status string    `json:"status"`
	At     time.Time `json:"at"`
}

type Job struct {
	ID                string    `json:"id"`
	Pair              string    `json:"pair"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
	Price             string    `json:"price,omitempty"`
	PriceDecimals     int       `json:"priceDecimals"`
	TxHash            string    `json:"txHash,omitempty"`
	BlockNumber       uint64    `json:"blockNumber,omitempty"`
	Error             string    `json:"error,omitempty"`
	Events            []Event   `json:"events"`
	SignedTransaction []byte    `json:"-"`
}

type storedJob struct {
	Job               Job    `json:"job"`
	SignedTransaction []byte `json:"signedTransaction,omitempty"`
}

func (j *Job) Transition(status string) {
	j.Status = status
	j.UpdatedAt = time.Now().UTC()
	j.Events = append(j.Events, Event{Status: status, At: j.UpdatedAt})
}

func jobKey(id string) string { return "oracle:job:" + id }

// EnqueueJob atomically admits a job and creates its stream entry. Duplicate IDs
// are rejected; failed Redis admission never returns a successful enqueue.
func (c *PriceCache) EnqueueJob(ctx context.Context, id, pair string, payload []byte) error {
	j := &Job{ID: id, Pair: pair, CreatedAt: time.Now().UTC(), PriceDecimals: 8, Events: []Event{}}
	j.Transition("queued")
	data, err := json.Marshal(storedJob{Job: *j})
	if err != nil {
		return err
	}
	return c.client.Eval(ctx, `
		if redis.call('EXISTS', KEYS[1]) == 1 then return redis.error_reply('duplicate job ID') end
		local streamType = redis.call('TYPE', KEYS[3]).ok
		local indexType = redis.call('TYPE', KEYS[2]).ok
		if streamType ~= 'none' and streamType ~= 'stream' then return redis.error_reply('invalid queue key type') end
		if indexType ~= 'none' and indexType ~= 'zset' then return redis.error_reply('invalid index key type') end
		redis.call('SET', KEYS[1], ARGV[1])
		redis.call('ZADD', KEYS[2], ARGV[2], ARGV[3])
		redis.call('ZREMRANGEBYRANK', KEYS[2], 0, -1001)
		redis.call('EXPIRE', KEYS[2], ARGV[5])
		return redis.call('XADD', KEYS[3], '*', 'job', ARGV[4])
	`, []string{jobKey(id), recentJobsKey, QueueStream}, data, j.CreatedAt.UnixMilli(), id, payload, int64(jobTTL/time.Second)).Err()
}

// SaveJob atomically persists the record and bounded recent index. Active jobs
// do not expire: a pending transaction must remain recoverable across downtime.
func (c *PriceCache) SaveJob(ctx context.Context, j *Job) error {
	data, err := json.Marshal(storedJob{Job: *j, SignedTransaction: j.SignedTransaction})
	if err != nil {
		return err
	}
	ttl := time.Duration(0)
	if j.Status == "confirmed" || j.Status == "cached" || j.Status == "failed" {
		ttl = jobTTL
	}
	_, err = c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		p.Set(ctx, jobKey(j.ID), data, ttl)
		p.ZAdd(ctx, recentJobsKey, redis.Z{Score: float64(j.CreatedAt.UnixMilli()), Member: j.ID})
		p.ZRemRangeByRank(ctx, recentJobsKey, 0, -1001)
		p.Expire(ctx, recentJobsKey, jobTTL)
		return nil
	})
	return err
}

func (c *PriceCache) GetJob(ctx context.Context, id string) (*Job, error) {
	data, err := c.client.Get(ctx, jobKey(id)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrJobNotFound
	}
	if err != nil {
		return nil, err
	}
	var stored storedJob
	if err := json.Unmarshal(data, &stored); err != nil {
		return nil, fmt.Errorf("decode job: %w", err)
	}
	stored.Job.SignedTransaction = stored.SignedTransaction
	return &stored.Job, nil
}

func (c *PriceCache) RecentJobs(ctx context.Context) ([]Job, error) {
	ids, err := c.client.ZRevRange(ctx, recentJobsKey, 0, 999).Result()
	if err != nil {
		return nil, err
	}
	jobs := make([]Job, 0, 30)
	for _, id := range ids {
		j, err := c.GetJob(ctx, id)
		if errors.Is(err, ErrJobNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *j)
		if len(jobs) == 30 {
			break
		}
	}
	return jobs, nil
}
