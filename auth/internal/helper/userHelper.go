package helper

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sachinggsingh/e-comm/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type JWT struct {
	Email string
	Uid   string
	jwt.RegisteredClaims
}

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrInGenerating = errors.New("error in generating token")
)

func GenerateToken(id string, email string) (string, string, error) {
	tokenclaims := &JWT{
		Email: email,
		Uid:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24 * time.Hour)),
		},
	}
	refreshTokenClaims := &JWT{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenclaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	tokenString, err := token.SignedString([]byte(config.GetEnv().APP_SECRET))
	refreshTokenString, err := refreshToken.SignedString([]byte(config.GetEnv().APP_SECRET))
	if err != nil {
		return "", "", ErrInGenerating
	}
	return tokenString, refreshTokenString, nil
}

func ValidateToken(token string) (*JWT, error) {
	tokenClaims := &JWT{}
	_, err := jwt.ParseWithClaims(token, tokenClaims, func(t *jwt.Token) (any, error) {
		return []byte(config.GetEnv().APP_SECRET), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	return tokenClaims, err
}

func Authorize(r *http.Request, userid string) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrInvalidToken
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", ErrInvalidToken
	}

	tokenString := authHeader[len(bearerPrefix):]

	claims := &JWT{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(config.GetEnv().APP_SECRET), nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}
	if claims.Uid != userid {
		return "", ErrInvalidToken
	}
	return claims.Uid, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
