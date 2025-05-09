package main

import (
	grpccache "api_service/internal/grpc_cache"
	"api_service/internal/handlers"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration",
			Buckets: []float64{0.1, 0.3, 1, 3, 5},
		},
		[]string{"path"},
	)

	cacheOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total cache operations",
		},
		[]string{"operation", "status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(cacheOperations)
}

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
	if err := http.ListenAndServe(":8051", mux); err != nil {
		log.Fatalf("Failed to start metrics/health server: %v", err)
	}
}

func main() {
	time.Sleep(time.Second * 3)
	// Запуск сервера метрик и health checks в отдельной горутине
	go startMetricsAndHealthServer()

	cache_service, err := grpccache.NewCacheClient("cache_service:50053")
	if err != nil {
		log.Fatal(err)
	}
	defer cache_service.Close()
	authHandler, err := handlers.NewUserAuthHandler("user_service:50051", cache_service)
	if err != nil {
		log.Fatal(err)
	}
	defer authHandler.Close()
	taskHandler, err := handlers.NewTaskServiceClient("task_service:50052", cache_service)
	if err != nil {
		log.Fatal(err)
	}
	defer taskHandler.Close()

	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		authHandler.RegisterRoutes(r)
	})
	r.Group(func(r chi.Router) {
		taskHandler.RegisterRoutes(r)
	})

	log.Println("Server starting on 3723:3723")
	if err := http.ListenAndServe(":3723", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
