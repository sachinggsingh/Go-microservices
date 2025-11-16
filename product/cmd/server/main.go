package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sachinggsingh/e-comm/internal/api"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/repository"
	"github.com/sachinggsingh/e-comm/internal/service"
)

func main() {
	env := config.GetEnv()
	database := db.NewDB()
	err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	server := api.NewServer(env, database)

	repo := repository.NewProductRepository(database)
	productService := service.NewProductService(repo)

	server.ProductRoutes(productService)

	go func() {
		if err := server.StartServer(); err != nil {
			fmt.Printf("Failed to start HTTP server: %v\n", err)
		}
	}()

	// Start gRPC Server
	go func() {
		if err := server.GrpcServer(); err != nil {
			fmt.Printf("Failed to start gRPC server: %v\n", err)
		}
		log.Printf("GRPC server running")
	}()

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down")

}
