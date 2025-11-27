package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cache "github.com/sachinggsingh/e-comm/internal/caches"
	"github.com/sachinggsingh/e-comm/internal/api"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/repository"
	"github.com/sachinggsingh/e-comm/internal/service"
)

func main() {
	env := config.GetEnv()

	// Initialize Redis for caching and rate limiting (single client)
	redisClient := config.NewRedisClient(env)
	redisCache := cache.NewRedisCache(redisClient)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
		_ = ctx // Context is used by cancel
	}()

	database := db.NewDB()
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	server := api.NewServer(env, database, redisCache)

	repo := repository.NewUserRepository(database)
	userService := service.NewUserService(repo)
	server.UserRoutes(userService)

	// Start HTTP server in a goroutine
	go func() {
		if err := server.StartServer(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start gRPC server in a goroutine
	// go func() {
	// 	api.StartGRPC()
	// }()

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down servers gracefully...")
}
