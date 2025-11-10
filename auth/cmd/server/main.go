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
	server := api.NewServer(env, database)

	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}
	defer database.Disconnect()

	repo := repository.NewUserRepository(database)
	userService := service.NewUserService(repo)
	server.UserRoutes(userService)

	server.StartServer()
}
