package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PORT       string
	MONGO_URL  string
	APP_SECRET string
}

func SetEnv() *Env {
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
		log.Fatalf("MONGODB_URI is not set")
	}
	appSecret := os.Getenv("APP_SECRET")
	if appSecret == "" {
		log.Fatalf("APP_SECRET is not set")
	}
	return &Env{
		PORT:       port,
		MONGO_URL:  mongoURL,
		APP_SECRET: appSecret,
	}
}
