package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(env *Env) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:         env.REDIS_HOST + ":" + env.REDIS_PORT,
		Password:     env.REDIS_PASSWORD,
		DB:           env.REDIS_DB,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")
	return client
}