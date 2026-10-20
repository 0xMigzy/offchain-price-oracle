package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/0xmigzy/capstone-worker/internal/cache"
	"github.com/0xmigzy/capstone-worker/internal/chain"
	"github.com/0xmigzy/capstone-worker/internal/queue"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// Worker is deliberately single-threaded: one signer must have one consumer.
type Worker struct {
	chain    *chain.Client
	cache    *cache.PriceCache
	consumer *queue.Consumer
}

func NewWorker(c *chain.Client, pc *cache.PriceCache, consumer *queue.Consumer) *Worker {
	return &Worker{chain: c, cache: pc, consumer: consumer}
}

func (w *Worker) Run(ctx context.Context) error { return w.consumer.Consume(ctx, w.processJob) }

func (w *Worker) fail(ctx context.Context, job *cache.Job, err error) error {
	job.Error = err.Error()
	job.Transition("failed")
	return w.cache.SaveJob(ctx, job)
}

func (w *Worker) processJob(ctx context.Context, request queue.Job) error {
	job, err := w.cache.GetJob(ctx, request.ID)
	if err != nil {
		return fmt.Errorf("load job %s: %w", request.ID, err)
	}
	if job.Pair != request.Pair {
		return fmt.Errorf("job pair mismatch: %s", request.ID)
	}
	switch job.Status {
	case "confirmed", "cached", "failed":
		return nil
	}
	if len(job.SignedTransaction) == 0 {
		job.Transition("processing")
		if err := w.cache.SaveJob(ctx, job); err != nil {
			return err
		}
		price, _, hit, err := w.cache.Get(ctx, job.Pair)
		if err != nil {
			return w.fail(ctx, job, errors.New("price cache unavailable"))
		}
		if hit {
			job.Price = price.String()
			job.Transition("cached")
			return w.cache.SaveJob(ctx, job)
		}
		// Intentionally simulated fixed prices, represented in raw 8-decimal units.
		switch job.Pair {
		case "ETH/USD":
			price = big.NewInt(3000 * 100000000)
		case "BTC/USD":
			price = big.NewInt(60000 * 100000000)
		default:
			return w.fail(ctx, job, fmt.Errorf("unsupported pair %q", job.Pair))
		}
		job.Price = price.String()
		rpcCtx, cancel := chain.RPCContext(ctx)
		tx, err := w.chain.PreparePrice(rpcCtx, job.Pair, price)
		cancel()
		if err != nil {
			return w.fail(ctx, job, errors.New("transaction preparation failed; check signer funding, authorization, and RPC availability"))
		}
		job.SignedTransaction, err = tx.MarshalBinary()
		if err != nil {
			return w.fail(ctx, job, err)
		}
		job.TxHash = tx.Hash().Hex()
		if err := w.cache.SaveJob(ctx, job); err != nil {
			return err
		}
	}

	var tx types.Transaction
	if err := tx.UnmarshalBinary(job.SignedTransaction); err != nil {
		return err
	}
	if job.Status != "submitted" {
		rpcCtx, cancel := chain.RPCContext(ctx)
		err := w.chain.Broadcast(rpcCtx, &tx)
		cancel()
		if err != nil {
			// A duplicate/already-mined transaction or lost send acknowledgement is
			// reconciled by hash. Unknown broadcasts stay recoverable, never failed.
			rpcCtx, cancel = chain.RPCContext(ctx)
			known, lookupErr := w.chain.TransactionKnown(rpcCtx, tx.Hash())
			cancel()
			if !known {
				return fmt.Errorf("broadcast %s: %w (lookup: %v)", job.TxHash, err, lookupErr)
			}
		}
		job.Transition("submitted")
		if err := w.cache.SaveJob(ctx, job); err != nil {
			return err
		}
	}

	for {
		rpcCtx, cancel := chain.RPCContext(ctx)
		receipt, err := w.chain.Receipt(rpcCtx, tx.Hash())
		cancel()
		if err == nil {
			job.BlockNumber = receipt.BlockNumber.Uint64()
			if receipt.Status != types.ReceiptStatusSuccessful {
				return w.fail(ctx, job, fmt.Errorf("transaction reverted"))
			}
			job.Transition("confirmed")
			if err := w.cache.SaveJob(ctx, job); err != nil {
				return err
			}
			price, _ := new(big.Int).SetString(job.Price, 10)
			if err := w.cache.Set(ctx, job.Pair, price, uint64(time.Now().Unix())); err != nil {
				log.Printf("confirmed job %s: cache update failed: %v", job.ID, err)
			}
			return nil
		}
		if !errors.Is(err, ethereum.NotFound) {
			return fmt.Errorf("receipt %s: %w", job.TxHash, err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(3 * time.Second):
		}
	}
}
