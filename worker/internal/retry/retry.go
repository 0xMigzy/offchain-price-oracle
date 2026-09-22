package retry

import (
	"context"
	"errors"
	"math"
	"time"
)

// Operation is a function that can be retried.
type Operation func(ctx context.Context) error

// Permanent wraps an error to indicate it should not be retried.
type Permanent struct{ Err error }

func (p *Permanent) Error() string { return p.Err.Error() }
func (p *Permanent) Unwrap() error { return p.Err }

// Config controls retry behaviour.
type Config struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// Default returns a sensible default config.
func Default() Config {
	return Config{
		MaxAttempts: 5,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    10 * time.Second,
	}
}

// Do executes the operation with retries.
func Do(ctx context.Context, cfg Config, op Operation) error {
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		err := op(ctx)
		if err == nil {
			return nil
		}

		var permanent *Permanent
		if errors.As(err, &permanent) {
			return permanent.Err // don't retry permanent errors
		}

		lastErr = err

		// Exponential backoff: BaseDelay * 2^attempt, capped at MaxDelay
		delay := cfg.BaseDelay * time.Duration(math.Pow(2, float64(attempt)))
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return lastErr
}
