package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xmigzy/capstone-worker/internal/chain"
	"github.com/0xmigzy/capstone-worker/internal/retry"
)

// OracleClient orchestrates oracle operations.
type OracleClient struct {
	chain *chain.Client
}

func NewOracleClient(c *chain.Client) *OracleClient {
	return &OracleClient{chain: c}
}

// PriceUpdate represents a price to submit.
type PriceUpdate struct {
	Pair  string
	Price *big.Int
}

// Submit submits a price update on chain.
func (o *OracleClient) Submit(ctx context.Context, update PriceUpdate) (string, error) {
	return o.chain.SubmitPrice(ctx, update.Pair, update.Price)
}

func (o *OracleClient) SubmitWithRetry(ctx context.Context, update PriceUpdate) (string, error) {
	var txHash string

	err := retry.Do(ctx, retry.Default(), func(ctx context.Context) error {
		h, err := o.chain.SubmitPrice(ctx, update.Pair, update.Price)
		if err != nil {
			return err
		}
		txHash = h
		return nil
	})

	return txHash, err
}

// GetLatest reads the latest price for a pair.
func (o *OracleClient) GetLatest(ctx context.Context, pair string) (*big.Int, uint64, error) {
	return o.chain.GetLatestPrice(ctx, pair)
}

// Close shuts down the client.
func (o *OracleClient) Close() error {
	// Add cleanup here when chain.Client gets a Close method
	return nil
}

// Verify does a sanity check that the chain client is working.
func (o *OracleClient) Verify(ctx context.Context) error {
	// Just call a known view function — if this fails, the connection is broken
	_, _, err := o.chain.GetLatestPrice(ctx, "ETH/USD")
	if err != nil {
		return fmt.Errorf("verify failed: %w", err)
	}
	return nil
}
