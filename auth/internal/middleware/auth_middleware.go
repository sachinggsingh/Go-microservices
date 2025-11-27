package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	cache "github.com/sachinggsingh/e-comm/internal/caches"
	"github.com/sachinggsingh/e-comm/internal/helper"
)

type AuthMiddleware struct {
	Redis *cache.RedisCache
}

func (a *AuthMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(token, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")

		tokenHash, _ := helper.HashToken(token)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.Redis.IsTokenBlacklisted(ctx, tokenHash); err != nil {
			http.Error(w, "token is expired or blacklisted", http.StatusUnauthorized)
			return
		}

		cachedUserID, err := a.Redis.IsTokenValid(ctx, tokenHash)
		if err == nil && cachedUserID != "" {
			r.Header.Set("user_id", cachedUserID)
			next.ServeHTTP(w, r)
			return
		}

		claims, err := helper.ValidateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		a.Redis.CacheValidToken(ctx, tokenHash, claims.Uid)

		r.Header.Set("user_id", claims.Uid)
		next.ServeHTTP(w, r)
	})
}
