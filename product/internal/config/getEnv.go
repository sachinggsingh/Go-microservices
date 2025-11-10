package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	MONGO_URL  string
	PORT       string
	APP_SECRET string
}

func GetEnv() *Env {
	load := godotenv.Load()
	if load != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not set")
	}
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Fatalf("MONGO_URL is not set")
	}

	app_sercert := os.Getenv("APP_SECRET")
	if app_sercert == "" {
		log.Fatalf("APP_SERCET is not set")
	}
	return &Env{
		PORT:       port,
		MONGO_URL:  mongoURL,
		APP_SECRET: app_sercert,
	}
}
