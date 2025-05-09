package main

import (
	"cache_service/internal/cache"
	"cache_service/internal/grpc/grpc_server"
	grpcclient "cache_service/internal/grpc_client"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var (
	cacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total cache hits",
		},
	)

	cacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total cache misses",
		},
	)

	cacheLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cache_operation_duration_seconds",
			Help:    "Cache operation duration",
			Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1},
		},
	)
)

func startMetricsAndHealthServer() {
	mux := http.NewServeMux()

	// Эндпоинт для метрик Prometheus
	mux.Handle("/metrics", promhttp.Handler())

	// Эндпоинт для health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Запуск сервера на порту 8051
	if err := http.ListenAndServe(":8052", mux); err != nil {
		log.Fatalf("Failed to start metrics/health server: %v", err)
	}
}

func init() {
	prometheus.MustRegister(cacheHits)
	prometheus.MustRegister(cacheMisses)
	prometheus.MustRegister(cacheLatency)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go startMetricsAndHealthServer()

	time.Sleep(3 * time.Second)
	rdb, err := cache.NewCache(&redis.Options{
		Addr:     "cache:6379",
		Password: "admin",
		DB:       0,
	})

	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	grpc_server.RegisterCacheServiceServer(s, &grpcclient.CacheServiceServer{Cch: rdb})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Server started at %v\n", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatal("Failed to serve: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server stopped")

}
