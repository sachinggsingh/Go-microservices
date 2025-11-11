package helper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sachinggsingh/e-comm/internal/config"
	"golang.org/x/crypto/bcrypt"
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

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrInGenerating = errors.New("error in generating token")
)

func GenerateToken(id string, email string) (string, string, error) {
	tokenclaims := &JWTAccessClaims{
		Email: email,
		Uid:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24 * time.Hour)),
		},
	}
	refreshTokenClaims := &JWTRefreshClaims{
		Uid: id,
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

func ValidateToken(tokenString string) (*JWTAccessClaims, error) {
	claims := &JWTAccessClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(config.GetEnv().APP_SECRET), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Uid == "" {
		return nil, errors.New("uid missing: likely refresh token used instead of access token")

	}
	return claims, nil
}
func Authorize(r *http.Request) (*JWTAccessClaims, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return nil, ErrInvalidToken
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, ErrInvalidToken
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	fmt.Println("Authenticated User:", claims.Uid)
	return claims, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
