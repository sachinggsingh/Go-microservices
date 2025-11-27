package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	cache "github.com/sachinggsingh/e-comm/internal/caches"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/helper"
	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userService *service.UserService
	redisCache  *cache.RedisCache
	env         *config.Env
}

func NewUserHandler(userService *service.UserService, redisCache *cache.RedisCache, env *config.Env) *UserHandler {
	return &UserHandler{
		userService: userService,
		redisCache:  redisCache,
		env:         env,
	}
}

// getClientIP extracts the client IP address with proper fallback handling
func (u *UserHandler) getClientIP(r *http.Request) string {
	// Try X-Forwarded-For header first (for proxies/load balancers)
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if ips := strings.Split(xForwardedFor, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	// Try X-Real-IP header as fallback
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return strings.TrimSpace(xRealIP)
	}
	// Fall back to RemoteAddr and strip port if needed
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := u.userService.RegisterUser(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	ip := u.getClientIP(r)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	attempts, _ := u.redisCache.IncrementLoginAttempt(ctx, ip, u.env.RATE_LIMIT_WINDOW_MINUTES)
	fmt.Println(attempts)
	if attempts > 5 {
		http.Error(w, "Too many attempts need to wait for 2 minutes", http.StatusTooManyRequests)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input data: "+err.Error(), http.StatusBadRequest)
		return
	}
	loggedInUser, err := u.userService.LoginUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.LoginResponse{
		ID:           loggedInUser.User_id,
		Email:        *loggedInUser.Email,
		Token:        *loggedInUser.Token,
		RefreshToken: *loggedInUser.Refresh_Token,
		User_id:      loggedInUser.User_id,
	}
	u.redisCache.ResetLoginAttempts(ctx, ip)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("user_id")
	if userIDFromRequest == "" {
		http.Error(w, "User id is required", http.StatusBadRequest)
		return
	}
	userID, err := helper.Authorize(r)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Authenticated User:", userID.Uid)
	user := &model.User{User_id: userIDFromRequest}
	result, err := u.userService.Profile(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}
