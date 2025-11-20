package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PORT                string
	MONGO_URL           string
	APP_SECRET          string
	PRODUCT_SERVICE_URL string
	STRIPE_SECRET_KEY   string
	STRIPE_SUCCESS_URL  string
	STRIPE_FAILURE_URL  string
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
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "localhost:9091" // Default to product service gRPC port
	}
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeSecretKey == "" {
		log.Fatalf("STRIPE_SECRET_KEY is not set")
	}
	stripeSuccessURL := os.Getenv("STRIPE_SUCCESS_URL")
	if stripeSuccessURL == "" {
		log.Fatalf("STRIPE_SUCCESS_URL is not set")
	}
	stripeFailureURL := os.Getenv("STRIPE_FAILURE_URL")
	if stripeFailureURL == "" {
		log.Fatalf("STRIPE_FAILURE_URL is not set")
	}
	return &Env{
		PORT:                port,
		MONGO_URL:           mongoURL,
		APP_SECRET:          appSecret,
		PRODUCT_SERVICE_URL: productServiceURL,
		STRIPE_SECRET_KEY:   stripeSecretKey,
		STRIPE_SUCCESS_URL:  stripeSuccessURL,
		STRIPE_FAILURE_URL:  stripeFailureURL,
	}
}
