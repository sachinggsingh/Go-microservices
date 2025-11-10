package middleware

import "github.com/sachinggsingh/e-comm/internal/config"

func GetUserIDFromToken(env *config.Env, userId string) (string, error) {
	return userId, nil
}
