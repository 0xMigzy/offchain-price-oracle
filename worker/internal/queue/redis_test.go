package queue

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/0xmigzy/capstone-worker/internal/cache"
	"github.com/redis/go-redis/v9"
)

// TEST_REDIS_URL must point at an isolated test Redis database. No FLUSH command
// is used; only this queue's well-known keys and explicitly created jobs change.
func testStore(t *testing.T) *cache.PriceCache {
	t.Helper()
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Skip("set TEST_REDIS_URL to an isolated Redis database")
	}
	store, err := cache.NewPriceCacheURL(url, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	keys := []string{cache.QueueStream, "oracle:jobs:recent", "oracle:job:queue-test"}
	if err := store.RedisClient().Del(context.Background(), keys...).Err(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.RedisClient().Del(context.Background(), keys...); store.Close() })
	return store
}

func TestPublishAtomicallyCreatesRecordAndStreamEntry(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()
	producer := NewProducer(store)
	job := Job{ID: "queue-test", Pair: "ETH/USD"}
	if err := producer.Publish(ctx, job); err != nil {
		t.Fatal(err)
	}
	record, err := store.GetJob(ctx, job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != "queued" || len(record.Events) != 1 || record.PriceDecimals != 8 {
		t.Fatalf("unexpected record: %+v", record)
	}
	if err := producer.Publish(ctx, job); err == nil {
		t.Fatal("duplicate ID accepted")
	}
	entries, err := store.RedisClient().XRange(ctx, cache.QueueStream, "-", "+").Result()
	if err != nil || len(entries) != 1 {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
	var decoded Job
	if err := json.Unmarshal([]byte(entries[0].Values["job"].(string)), &decoded); err != nil || decoded != job {
		t.Fatalf("payload=%+v err=%v", decoded, err)
	}
}

func TestAdmissionFailureDoesNotCreateJob(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()
	if err := store.RedisClient().Set(ctx, cache.QueueStream, "wrong-type", 0).Err(); err != nil {
		t.Fatal(err)
	}
	if err := NewProducer(store).Publish(ctx, Job{ID: "queue-test", Pair: "ETH/USD"}); err == nil {
		t.Fatal("invalid stream accepted")
	}
	if _, err := store.GetJob(ctx, "queue-test"); !errors.Is(err, cache.ErrJobNotFound) {
		t.Fatalf("failed admission left job: %v", err)
	}
}

func TestConsumerRecoversPendingEntryAndDeletesOnlyAfterSuccess(t *testing.T) {
	store := testStore(t)
	ctx := context.Background()
	if err := NewProducer(store).Publish(ctx, Job{ID: "queue-test", Pair: "ETH/USD"}); err != nil {
		t.Fatal(err)
	}
	consumer := NewConsumer(store.RedisClient())
	if err := consumer.initialize(ctx); err != nil {
		t.Fatal(err)
	}
	// Deliver once without ACK, simulating a process dying while working.
	if _, err := store.RedisClient().XReadGroup(ctx, &redis.XReadGroupArgs{Group: group, Consumer: consumerName, Streams: []string{cache.QueueStream, ">"}, Count: 1}).Result(); err != nil {
		t.Fatal(err)
	}
	failedCtx, cancel := context.WithCancel(ctx)
	_ = consumer.Consume(failedCtx, func(context.Context, Job) error { cancel(); return errors.New("storage unavailable") })
	pending, err := store.RedisClient().XPending(ctx, cache.QueueStream, group).Result()
	if err != nil || pending.Count != 1 {
		t.Fatalf("pending=%+v err=%v", pending, err)
	}
	restarted := NewConsumer(store.RedisClient())
	runCtx, stop := context.WithTimeout(ctx, 1500*time.Millisecond)
	defer stop()
	calls := 0
	_ = restarted.Consume(runCtx, func(ctx context.Context, job Job) error {
		calls++
		record, err := store.GetJob(ctx, job.ID)
		if err != nil {
			return err
		}
		record.Error = "simulated terminal failure"
		record.Transition("failed")
		return store.SaveJob(ctx, record)
	})
	if calls != 1 {
		t.Fatalf("recovered handler calls=%d", calls)
	}
	if size, err := store.RedisClient().XLen(ctx, cache.QueueStream).Result(); err != nil || size != 0 {
		t.Fatalf("remaining entries=%d err=%v", size, err)
	}
	pending, err = store.RedisClient().XPending(ctx, cache.QueueStream, group).Result()
	if err != nil || pending.Count != 0 {
		t.Fatalf("pending=%+v err=%v", pending, err)
	}
	record, err := store.GetJob(ctx, "queue-test")
	if err != nil || record.Status != "failed" {
		t.Fatalf("terminal job=%+v err=%v", record, err)
	}
}
