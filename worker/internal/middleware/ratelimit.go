package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter *rate.Limiter
	seen    time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rps      rate.Limit
	burst    int
}

func NewRateLimiter(rps float64, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*ipLimiter),
		rps:      rate.Limit(rps),
		burst:    burst,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !rl.allow(ip) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	l, ok := rl.limiters[ip]
	if !ok {
		l = &ipLimiter{limiter: rate.NewLimiter(rl.rps, rl.burst)}
		rl.limiters[ip] = l
	}
	l.seen = time.Now()
	return l.limiter.Allow()
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(5 * time.Minute)
		rl.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for ip, l := range rl.limiters {
			if l.seen.Before(cutoff) {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func clientIP(r *http.Request) string {
	// Only trust the socket peer; public clients can forge forwarded headers.
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}
