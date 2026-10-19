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
	"github.com/0xmigzy/capstone-worker/internal/client"
	"github.com/0xmigzy/capstone-worker/internal/queue"
	"github.com/0xmigzy/capstone-worker/internal/worker"
)

const (
	kafkaTopic   = "oracle-jobs"
	kafkaGroupID = "oracle-worker"
	cacheTTL     = 10 * time.Second
	httpAddr     = "0.0.0.0:8080"
	httpShutdown = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rpcURL := mustEnv("RPC_URL")
	contractAddr := mustEnv("CONTRACT_ADDRESS")
	privateKey := strings.TrimPrefix(mustEnv("PRIVATE_KEY"), "0x")
	redisAddr := mustEnv("REDIS_ADDR")
	brokers := strings.Split(mustEnv("KAFKA_BROKERS"), ",")

	chainClient, err := chain.NewClient(ctx, rpcURL, contractAddr, privateKey)
	if err != nil {
		log.Fatalf("create chain client: %v", err)
	}

	priceCache := cache.NewPriceCache(redisAddr, cacheTTL)
	oracle := client.NewOracleClient(chainClient, priceCache)
	consumer := queue.NewConsumer(brokers, kafkaTopic, kafkaGroupID)

	w := worker.NewWorker(oracle, priceCache, consumer)

	producer := queue.NewProducer(brokers, kafkaTopic)
	s := api.NewServer(producer)
	mux := healthMux()
	s.Routes(mux)

	ln, err := net.Listen("tcp", httpAddr)
	if err != nil {
		log.Fatalf("http listen %s: %v", httpAddr, err)
	}
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	srvErr := make(chan error, 1)
	go func() {
		if err := srv.Serve(ln); !errors.Is(err, http.ErrServerClosed) {
			srvErr <- err
			stop() // take the worker down with the server
		}
	}()
	log.Printf("http listening on %s", httpAddr)

	log.Printf("worker started: brokers=%v topic=%s group=%s", brokers, kafkaTopic, kafkaGroupID)
	err = w.Run(ctx)
	if cerr := consumer.Close(); cerr != nil {
		log.Printf("close consumer: %v", cerr)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdown)
	defer cancel()
	if serr := srv.Shutdown(shutdownCtx); serr != nil {
		log.Printf("http shutdown: %v", serr)
	}
	if perr := producer.Close(); perr != nil {
		log.Printf("close producer: %v", perr)
	}
	select {
	case serr := <-srvErr:
		log.Fatalf("http server: %v", serr)
	default:
	}
	if err != nil && ctx.Err() == nil {
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
