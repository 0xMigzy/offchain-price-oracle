package worker

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/0xmigzy/capstone-worker/internal/cache"
	"github.com/0xmigzy/capstone-worker/internal/client"
	"github.com/0xmigzy/capstone-worker/internal/queue"
)

// Worker pulls jobs from Kafka, checks the cache, and submits prices on chain.
type Worker struct {
	oracle   *client.OracleClient
	cache    *cache.PriceCache
	consumer *queue.Consumer
}

// NewWorker creates a Worker.
func NewWorker(oracle *client.OracleClient, pc *cache.PriceCache, consumer *queue.Consumer) *Worker {
	return &Worker{oracle: oracle, cache: pc, consumer: consumer}
}

// Run pulls jobs from Kafka and processes each one until ctx is cancelled.
func (w *Worker) Run(ctx context.Context) error {
	return w.consumer.Consume(ctx, w.processJob)
}

// simulatePrice returns a fixed fake price for testing.
func simulatePrice(pair string) float64 {
	switch pair {
	case "ETH/USD":
		return 3000.00
	case "BTC/USD":
		return 60000.00
	default:
		return 1.00
	}
}

func (w *Worker) processJob(ctx context.Context, job queue.Job) error {
	// Check cache
	_, ts, hit, _ := w.cache.Get(ctx, job.Pair)
	if hit && time.Now().Unix()-int64(ts) < 60 {
		// Cached, recent — skip
		return nil
	}

	// Simulate fetching price (or call real source)
	fetched := simulatePrice(job.Pair)
	priceWei := big.NewInt(int64(fetched * 1e8))

	// Submit on chain
	txHash, err := w.oracle.SubmitWithRetry(ctx, client.PriceUpdate{
		Pair:  job.Pair,
		Price: priceWei,
	})
	if err != nil {
		return err
	}

	// Update cache
	_ = w.cache.Set(ctx, job.Pair, priceWei, uint64(time.Now().Unix()))

	log.Printf("submitted: pair=%s tx=%s", job.Pair, txHash)
	return nil
}
