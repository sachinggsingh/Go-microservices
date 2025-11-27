package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (rc *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rc.client.Set(ctx, key, data, ttl).Err()
}

func (rc *RedisCache) Get(ctx context.Context, key string, value any) error {
	data, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), value)
}

func (rc *RedisCache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

func (rc *RedisCache) SetString(ctx context.Context, key string, value string, ttl time.Duration) error {
	return rc.client.Set(ctx, key, value, ttl).Err()
}

func (rc *RedisCache) GetString(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}
