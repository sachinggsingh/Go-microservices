package cache

import (
	"context"
	"errors"
	"time"
)

func (rc *RedisCache) BlacklistToken(ctx context.Context, tokenHash string) error {
	key := "token:blacklist:" + tokenHash
	return rc.SetString(ctx, key, "1", time.Hour*24)
}

func (rc *RedisCache) IsTokenBlacklisted(ctx context.Context, tokenHash string) error {
	key := "token:blacklist:" + tokenHash
	val, err := rc.GetString(ctx, key)
	if err == nil && val != "" {
		return errors.New("token blacklisted")
	}
	return nil
}
