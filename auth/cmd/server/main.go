package main

import (
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
	server := api.NewServer(env, database)

	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer database.Disconnect()

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
	go func() {
		api.StartGRPC()
	}()

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down servers")
}
