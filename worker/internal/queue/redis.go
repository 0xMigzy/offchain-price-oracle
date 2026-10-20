package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/0xmigzy/capstone-worker/internal/cache"
	"github.com/redis/go-redis/v9"
)

const group = "oracle-worker"

// A stable consumer identity recovers its pending entries after a restart.
// Exactly one worker replica is required, also to serialize the signing nonce.
const consumerName = "single-worker"

type Job struct {
	ID   string `json:"id"`
	Pair string `json:"pair"`
}

type Producer struct{ store *cache.PriceCache }

func NewProducer(store *cache.PriceCache) *Producer { return &Producer{store: store} }

func (p *Producer) Publish(ctx context.Context, job Job) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}
	return p.store.EnqueueJob(ctx, job.ID, job.Pair, body)
}

func (p *Producer) Check(ctx context.Context) error {
	groups, err := p.store.RedisClient().XInfoGroups(ctx, cache.QueueStream).Result()
	if err != nil {
		return err
	}
	for _, info := range groups {
		if info.Name == group {
			return nil
		}
	}
	return fmt.Errorf("queue consumer group unavailable")
}

type Consumer struct{ client *redis.Client }

func NewConsumer(client *redis.Client) *Consumer { return &Consumer{client: client} }

func (c *Consumer) initialize(ctx context.Context) error {
	err := c.client.XGroupCreateMkStream(ctx, cache.QueueStream, group, "0").Err()
	if err != nil && !strings.HasPrefix(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}

// Consume drains pending entries before reading new work. Failed handlers retain
// the stream entry and pending ownership; terminal persistence precedes XACK.
func (c *Consumer) Consume(ctx context.Context, handler func(context.Context, Job) error) error {
	if err := c.initialize(ctx); err != nil {
		return err
	}
	pending := true
	for {
		id := ">"
		block := time.Second
		if pending {
			id = "0"
			block = -1
		}
		streams, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group: group, Consumer: consumerName, Streams: []string{cache.QueueStream, id}, Count: 1, Block: block,
		}).Result()
		if errors.Is(err, redis.Nil) || (err == nil && (len(streams) == 0 || len(streams[0].Messages) == 0)) {
			pending = false
			if ctx.Err() != nil {
				return ctx.Err()
			}
			continue
		}
		if err != nil {
			return err
		}
		for _, stream := range streams {
			for _, message := range stream.Messages {
				body, ok := message.Values["job"].(string)
				if !ok {
					return fmt.Errorf("invalid queue entry %s retained: missing job payload", message.ID)
				}
				var job Job
				if err := json.Unmarshal([]byte(body), &job); err != nil || job.ID == "" || job.Pair == "" {
					return fmt.Errorf("invalid queue entry %s retained", message.ID)
				}
				for {
					if err := handler(ctx, job); err == nil {
						break
					} else {
						log.Printf("job %s retained for retry: %v", job.ID, err)
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case <-time.After(5 * time.Second):
					}
				}
				// A single group owns this stream; delete acknowledged entries so
				// trimming can never discard unprocessed or pending work.
				_, err := c.client.TxPipelined(ctx, func(p redis.Pipeliner) error {
					p.XAck(ctx, cache.QueueStream, group, message.ID)
					p.XDel(ctx, cache.QueueStream, message.ID)
					return nil
				})
				if err != nil {
					return err
				}
			}
		}
	}
}
