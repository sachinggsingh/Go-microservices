package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/errors"
)

type JWTAccessClaims struct {
	Email string `json:"email"`
	Uid   string `json:"uid"`
	jwt.RegisteredClaims
}

type JWTRefreshClaims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

type contextKey string

const UserIDKey contextKey = "uid"

// GetTheToken extracts and validates the JWT token.
func GetTheToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.ErrNoTokenProvided
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", errors.ErrInvalidToken
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
func ValidateToken(tokenString string) (*JWTAccessClaims, error) {
	claims := &JWTAccessClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.SetEnv().APP_SECRET), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.ErrInvalidToken
	}
	if claims.Uid == "" {
		return nil, errors.DifferentTokenUsed

	}
	return claims, nil
}

// Middleware: Inject user_id into request context
func GetUserIdFromToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString, err := GetTheToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// ✅ Validate ACCESS TOKEN properly
		claims, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		fmt.Println("Authenticated User:", claims)

		// Now we know it's a valid ACCESS token
		userID := claims.Uid
		if userID == "" {
			http.Error(w, "Invalid access token: uid missing", http.StatusUnauthorized)
			return
		}
		fmt.Println("Authenticated User:", userID)

		// Inject UID in context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
