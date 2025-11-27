package cache

import (
	"context"
	"time"
)

func (rc *RedisCache) CacheValidToken(ctx context.Context, tokenHash string, userID string) error {
	key := "token:valid:" + tokenHash
	return rc.SetString(ctx, key, userID, time.Hour)
}

func (rc *RedisCache) IsTokenValid(ctx context.Context, tokenHash string) (string, error) {
	key := "token:valid:" + tokenHash
	return rc.GetString(ctx, key)
}
