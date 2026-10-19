// internal/cache/redis_test.go
package cache

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestPriceCache_GetSet(t *testing.T) {
	// Skip if Redis not available
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := client.Ping(context.Background()).Err(); err != nil {
		t.Skip("Redis not available")
	}

	cache := NewPriceCache("localhost:6379", 1*time.Minute)
	ctx := context.Background()

	// Set
	err := cache.Set(ctx, "ETH/USD", big.NewInt(3500e8), 1234567890)
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	// Get
	price, ts, hit, err := cache.Get(ctx, "ETH/USD")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if !hit {
		t.Fatal("expected cache hit")
	}
	if price.Cmp(big.NewInt(3500e8)) != 0 {
		t.Errorf("price mismatch: got %s, want 3500e8", price.String())
	}
	if ts != 1234567890 {
		t.Errorf("timestamp mismatch: got %d, want 1234567890", ts)
	}
}
