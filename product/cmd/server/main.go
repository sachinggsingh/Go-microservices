package main

import (
	"log"

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

	server.StartServer()
}
