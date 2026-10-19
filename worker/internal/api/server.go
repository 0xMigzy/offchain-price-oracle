// Package api is the submitter HTTP API: it turns oracle requests into Kafka jobs.
package api

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/0xmigzy/capstone-worker/internal/middleware"
	"github.com/0xmigzy/capstone-worker/internal/queue"
)

const maxBodyBytes = 1 << 10

// Server publishes accepted oracle requests to the jobs topic.
type Server struct {
	producer *queue.Producer
}

func NewServer(producer *queue.Producer) *Server {
	return &Server{producer: producer}
}

// Routes registers the API handlers on mux.
func (s *Server) Routes(mux *http.ServeMux) {
	rateLimiter := middleware.NewRateLimiter(10, 20) // 10 rps, burst 20
	mux.Handle("/oracle/request", rateLimiter.Middleware(http.HandlerFunc(s.handleOracleRequest)))
}

type oracleRequest struct {
	Pair string `json:"pair"`
}

type oracleResponse struct {
	ID     string `json:"id"`
	Pair   string `json:"pair"`
	Status string `json:"status"`
}

func (s *Server) handleOracleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req oracleRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	pair := strings.TrimSpace(req.Pair)
	if pair == "" {
		http.Error(w, "pair is required", http.StatusBadRequest)
		return
	}

	job := queue.Job{ID: rand.Text(), Pair: pair}
	if err := s.producer.Publish(r.Context(), job); err != nil {
		log.Printf("api: publish job %s: %v", job.ID, err)
		http.Error(w, "failed to enqueue job", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(oracleResponse{ID: job.ID, Pair: job.Pair, Status: "queued"})
}
