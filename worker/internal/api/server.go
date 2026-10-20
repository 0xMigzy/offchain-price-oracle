// Package api exposes oracle submissions and truthful demo status reads.
package api

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/0xmigzy/capstone-worker/internal/cache"
	"github.com/0xmigzy/capstone-worker/internal/chain"
	"github.com/0xmigzy/capstone-worker/internal/middleware"
	"github.com/0xmigzy/capstone-worker/internal/queue"
)

const maxBodyBytes = 1 << 10

type Server struct {
	producer *queue.Producer
	cache    *cache.PriceCache
	chain    *chain.Client
	token    string
	contract string
}

func NewServer(producer *queue.Producer, pc *cache.PriceCache, c *chain.Client, token, contract string) *Server {
	return &Server{producer: producer, cache: pc, chain: c, token: token, contract: contract}
}

func (s *Server) Routes(mux *http.ServeMux) {
	rateLimiter := middleware.NewRateLimiter(10, 20)
	mux.Handle("POST /oracle/request", rateLimiter.Middleware(http.HandlerFunc(s.handleOracleRequest)))
	mux.HandleFunc("GET /api/system", s.handleSystem)
	mux.HandleFunc("GET /api/jobs", s.handleJobs)
	mux.HandleFunc("GET /api/jobs/{id}", s.handleJob)
	mux.HandleFunc("GET /api/prices", s.handlePrices)
}

func validPair(pair string) bool { return pair == "ETH/USD" || pair == "BTC/USD" }

func jsonResponse(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (s *Server) handleOracleRequest(w http.ResponseWriter, r *http.Request) {
	if s.token != "" {
		provided := sha256.Sum256([]byte(r.Header.Get("Authorization")))
		expected := sha256.Sum256([]byte("Bearer " + s.token))
		if subtle.ConstantTimeCompare(provided[:], expected[:]) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
	var req struct {
		Pair string `json:"pair"`
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	pair := strings.TrimSpace(req.Pair)
	if !validPair(pair) {
		http.Error(w, "supported pairs: ETH/USD, BTC/USD", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	job := queue.Job{ID: rand.Text(), Pair: pair}
	if err := s.producer.Publish(ctx, job); err != nil {
		log.Printf("enqueue job %s: %v", job.ID, err)
		http.Error(w, "failed to enqueue job", http.StatusServiceUnavailable)
		return
	}
	jsonResponse(w, http.StatusAccepted, map[string]string{"id": job.ID, "pair": job.Pair, "status": "queued"})
}

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	jobs, err := s.cache.RecentJobs(ctx)
	if err != nil {
		http.Error(w, "job storage unavailable", http.StatusServiceUnavailable)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"jobs": jobs})
}

func (s *Server) handleJob(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	job, err := s.cache.GetJob(ctx, r.PathValue("id"))
	if errors.Is(err, cache.ErrJobNotFound) {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "job storage unavailable", http.StatusServiceUnavailable)
		return
	}
	jsonResponse(w, http.StatusOK, job)
}

func (s *Server) handlePrices(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")
	if !validPair(pair) {
		http.Error(w, "supported pairs: ETH/USD, BTC/USD", http.StatusBadRequest)
		return
	}
	ctx, cancel := chain.RPCContext(r.Context())
	defer cancel()
	price, timestamp, err := s.chain.GetLatestPrice(ctx, pair)
	if err != nil {
		http.Error(w, "chain read unavailable", http.StatusServiceUnavailable)
		return
	}
	jsonResponse(w, http.StatusOK, map[string]any{"pair": pair, "price": price.String(), "timestamp": timestamp, "priceDecimals": 8})
}

func (s *Server) handleSystem(w http.ResponseWriter, r *http.Request) {
	type checkResult struct{ name, status string }
	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()
	results := make(chan checkResult, 3)
	checks := map[string]func(context.Context) error{
		"rpc":   func(ctx context.Context) error { _, _, err := s.chain.GetLatestPrice(ctx, "ETH/USD"); return err },
		"redis": s.cache.Ping,
		"queue": s.producer.Check,
	}
	for name, check := range checks {
		go func() {
			status := "healthy"
			if check(ctx) != nil {
				status = "unavailable"
			}
			results <- checkResult{name, status}
		}()
	}
	services := map[string]string{}
	for range 3 {
		result := <-results
		services[result.name] = result.status
	}
	explorer := ""
	if s.chain.NetworkID() == 11155111 {
		explorer = "https://sepolia.etherscan.io"
	}
	jsonResponse(w, http.StatusOK, map[string]any{
		"chainId": s.chain.NetworkID(), "contractAddress": s.contract, "explorerBaseUrl": explorer,
		"simulated": true, "queueBackend": "redis-streams", "services": services, "pairs": []string{"ETH/USD", "BTC/USD"},
	})
}
