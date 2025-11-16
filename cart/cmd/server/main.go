package main

import (
	"log"

	restapi "github.com/sachinggsingh/e-comm/internal/api"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/pkg"
	"github.com/sachinggsingh/e-comm/internal/repository"
	"github.com/sachinggsingh/e-comm/internal/service"
)

func main() {
	env := config.SetEnv()
	database := db.NewDatabase()
	err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	// Initialize product gRPC client
	productClient, err := pkg.NewProductClient(env.PRODUCT_SERVICE_URL)
	if err != nil {
		log.Fatalf("Failed to initialize product gRPC client: %v", err)
	}
	defer productClient.Close()

	server := restapi.NewServer(env, database)
	repo := repository.NewCartRepository(database)
	cartService := service.NewCartService(repo, productClient)
	server.CartRoute(cartService)
	server.StartServer()
}
