package config

import (
	"github.com/daddydemir/crypto/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	host := config.Get("REDIS_HOST")
	pass := config.Get("REDIS_PASS")

	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       0,
	})
}
