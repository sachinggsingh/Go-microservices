package main

import (
	"log"

	restapi "github.com/sachinggsingh/e-comm/internal/api"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
)

func main() {
	env := config.SetEnv()
	database := db.NewDatabase()
	err := database.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	server := restapi.NewServer(env, database)
	server.StartServer()
}
