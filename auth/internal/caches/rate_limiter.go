package cache

import (
	"context"
	"log"
	"time"
)

// IncrementLoginAttempt increments the login attempt counter for a given IP
// and sets expiration if it's the first attempt
func (rc *RedisCache) IncrementLoginAttempt(ctx context.Context, ip string, windowMinutes int) (int, error) {
	key := "ratelimit:login:" + ip
	count, err := rc.client.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Error incrementing rate limit for IP %s: %v", ip, err)
		return 0, err
	}
	if count == 1 {
		// Set expiration only on first increment
		rc.client.Expire(ctx, key, time.Duration(windowMinutes)*time.Minute)
	}
	return int(count), nil
}

// ResetLoginAttempts resets the login attempt counter for a given IP
func (rc *RedisCache) ResetLoginAttempts(ctx context.Context, ip string) error {
	key := "ratelimit:login:" + ip
	if err := rc.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Error resetting rate limit for IP %s: %v", ip, err)
		return err
	}
	return nil
}

// IncrementRegisterAttempt increments the register attempt counter for a given IP
func (rc *RedisCache) IncrementRegisterAttempt(ctx context.Context, ip string, windowMinutes int) (int, error) {
	key := "ratelimit:register:" + ip
	count, err := rc.client.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("Error incrementing register rate limit for IP %s: %v", ip, err)
		return 0, err
	}
	if count == 1 {
		rc.client.Expire(ctx, key, time.Duration(windowMinutes)*time.Minute)
	}
	return int(count), nil
}

// ResetRegisterAttempts resets the register attempt counter for a given IP
func (rc *RedisCache) ResetRegisterAttempts(ctx context.Context, ip string) error {
	key := "ratelimit:register:" + ip
	if err := rc.client.Del(ctx, key).Err(); err != nil {
		log.Printf("Error resetting register rate limit for IP %s: %v", ip, err)
		return err
	}
	return nil
}
