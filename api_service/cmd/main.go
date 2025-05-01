package main

import (
	grpccache "api_service/internal/grpc_cache"
	"api_service/internal/handlers"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	time.Sleep(time.Second * 3)
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
