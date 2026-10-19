package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

// PriceCache caches price reads.
type PriceCache struct {
	client *redis.Client
	ttl    time.Duration
}

// CachedPrice is what we store in Redis.
type CachedPrice struct {
	Price     string `json:"price"` // big.Int serialised as string
	Timestamp uint64 `json:"timestamp"`
}

func NewPriceCache(addr string, ttl time.Duration) *PriceCache {
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &PriceCache{client: client, ttl: ttl}
}

func (c *PriceCache) Get(ctx context.Context, pair string) (*big.Int, uint64, bool, error) {
	val, err := c.client.Get(ctx, key(pair)).Bytes()
	if err == redis.Nil {
		return nil, 0, false, nil // cache miss
	}
	if err != nil {
		return nil, 0, false, fmt.Errorf("redis get: %w", err)
	}

	var cached CachedPrice
	if err := json.Unmarshal(val, &cached); err != nil {
		return nil, 0, false, fmt.Errorf("unmarshal: %w", err)
	}

	price, ok := new(big.Int).SetString(cached.Price, 10)
	if !ok {
		return nil, 0, false, fmt.Errorf("parse price: %s", cached.Price)
	}

	return price, cached.Timestamp, true, nil
}

func (c *PriceCache) Set(ctx context.Context, pair string, price *big.Int, ts uint64) error {
	cached := CachedPrice{
		Price:     price.String(),
		Timestamp: ts,
	}
	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := c.client.Set(ctx, key(pair), data, c.ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}
	return nil
}

func key(pair string) string { return "price:" + pair }
