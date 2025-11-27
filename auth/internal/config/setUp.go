package config

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	PORT           string
	MONGO_URL      string
	APP_SECRET     string
	REDIS_HOST     string
	REDIS_PORT     string
	REDIS_PASSWORD string
	REDIS_DB              int
	RATE_LIMIT_MAX_ATTEMPTS int
	RATE_LIMIT_WINDOW_MINUTES int
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

	host := os.Getenv("REDIS_HOST")

	if host == "redis" {
		_, err := net.LookupHost("redis")
		if err != nil {
			host = "localhost"
		}
	}

	redis_port := os.Getenv("REDIS_PORT")
	if redis_port == "" {
		log.Fatalf("REDIS_PORT is not set")
	}
	redis_password := os.Getenv("REDIS_PASSWORD")
	if redis_password == "" {
		log.Fatalf("REDIS_PASSWORD is not set")
	}
	redis_db_str := os.Getenv("REDIS_DB")
	if redis_db_str == "" {
		log.Fatalf("REDIS_DB is not set")
	}
	redis_db, err := strconv.Atoi(redis_db_str)
	if err != nil {
		log.Fatalf("REDIS_DB must be a valid integer: %v", err)
	}

	// Rate limiting configuration with defaults
	rate_limit_attempts_str := os.Getenv("RATE_LIMIT_MAX_ATTEMPTS")
	if rate_limit_attempts_str == "" {
		rate_limit_attempts_str = "5"
	}
	rate_limit_attempts, err := strconv.Atoi(rate_limit_attempts_str)
	if err != nil {
		log.Fatalf("RATE_LIMIT_MAX_ATTEMPTS must be a valid integer: %v", err)
	}

	rate_limit_window_str := os.Getenv("RATE_LIMIT_WINDOW_MINUTES")
	if rate_limit_window_str == "" {
		rate_limit_window_str = "5"
	}
	rate_limit_window, err := strconv.Atoi(rate_limit_window_str)
	if err != nil {
		log.Fatalf("RATE_LIMIT_WINDOW_MINUTES must be a valid integer: %v", err)
	}

	return &Env{
		PORT:                      port,
		MONGO_URL:                 mongoURL,
		APP_SECRET:                app_sercert,
		REDIS_HOST:                host,
		REDIS_PORT:                redis_port,
		REDIS_PASSWORD:            redis_password,
		REDIS_DB:                  redis_db,
		RATE_LIMIT_MAX_ATTEMPTS:   rate_limit_attempts,
		RATE_LIMIT_WINDOW_MINUTES: rate_limit_window,
	}
}
