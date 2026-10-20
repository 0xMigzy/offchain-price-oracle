package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/0xmigzy/capstone-worker/internal/api"
	"github.com/0xmigzy/capstone-worker/internal/cache"
	"github.com/0xmigzy/capstone-worker/internal/chain"
	"github.com/0xmigzy/capstone-worker/internal/queue"
	"github.com/0xmigzy/capstone-worker/internal/worker"
)

const (
	cacheTTL     = 10 * time.Second
	httpShutdown = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rpcURL := mustEnv("RPC_URL")
	contractAddr := mustEnv("CONTRACT_ADDRESS")
	privateKey := strings.TrimPrefix(mustEnv("PRIVATE_KEY"), "0x")
	redisURL := os.Getenv("REDIS_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisURL == "" && redisAddr == "" {
		log.Fatal("REDIS_URL or REDIS_ADDR is required")
	}
	token := os.Getenv("BACKEND_API_TOKEN")
	if token == "" && os.Getenv("ALLOW_UNAUTHENTICATED_LOCAL") != "true" {
		log.Fatal("BACKEND_API_TOKEN is required (local-only opt out: ALLOW_UNAUTHENTICATED_LOCAL=true)")
	}
	if os.Getenv("RAILWAY_ENVIRONMENT_ID") != "" && token == "" {
		log.Fatal("hosted deployment requires BACKEND_API_TOKEN")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	httpAddr := net.JoinHostPort("0.0.0.0", port)

	startupCtx, startupCancel := chain.RPCContext(ctx)
	chainClient, err := chain.NewClient(startupCtx, rpcURL, contractAddr, privateKey)
	startupCancel()
	if err != nil {
		log.Fatalf("create chain client: %v", err)
	}

	defer chainClient.Close()
	if os.Getenv("RAILWAY_ENVIRONMENT_ID") != "" && chainClient.NetworkID() != 11155111 {
		log.Fatal("hosted demo requires Sepolia chain ID 11155111")
	}
	var priceCache *cache.PriceCache
	if redisURL != "" {
		priceCache, err = cache.NewPriceCacheURL(redisURL, cacheTTL)
		if err != nil {
			log.Fatalf("create Redis client: %v", err)
		}
	} else {
		priceCache = cache.NewPriceCache(redisAddr, cacheTTL)
	}
	defer priceCache.Close()
	consumer := queue.NewConsumer(priceCache.RedisClient())

	w := worker.NewWorker(chainClient, priceCache, consumer)

	producer := queue.NewProducer(priceCache)
	s := api.NewServer(producer, priceCache, chainClient, token, contractAddr)
	mux := healthMux()
	s.Routes(mux)

	ln, err := net.Listen("tcp", httpAddr)
	if err != nil {
		log.Fatalf("http listen %s: %v", httpAddr, err)
	}
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 20 * time.Second, WriteTimeout: 25 * time.Second, IdleTimeout: 60 * time.Second}
	srvErr := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
			stop() // take the worker down with the server
		}
	}()
	log.Printf("http listening on %s", httpAddr)

	log.Printf("worker started: queue=redis-streams stream=%s", cache.QueueStream)
	err = w.Run(ctx)
	stop()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdown)
	defer cancel()
	if serr := srv.Shutdown(shutdownCtx); serr != nil {
		log.Printf("http shutdown: %v", serr)
	}
	select {
	case serr := <-srvErr:
		log.Fatalf("http server: %v", serr)
	default:
	}
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("worker: %v", err)
	}
	log.Printf("worker stopped")
}

func healthMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env var %s", key)
	}
	return v
}
