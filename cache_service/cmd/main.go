package main

import (
	"cache_service/internal/cache"
	"cache_service/internal/grpc/grpc_server"
	grpcclient "cache_service/internal/grpc_client"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func main() {
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
