package cache

import (
	"context"
	"time"
)

type UserSession struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func (rc *RedisCache) SaveUserSession(ctx context.Context, userID string, userSession UserSession) error {
	key := "user:session:" + userID
	return rc.Set(ctx, key, userSession, 30*time.Minute)
}

func (rc *RedisCache) GetSavedUserSession(ctx context.Context, userID string) (*UserSession, error) {
	key := "user:session:" + userID

	var session UserSession
	err := rc.Get(ctx, key, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}
