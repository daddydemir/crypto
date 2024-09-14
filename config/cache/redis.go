package cache

import (
	"context"
	"github.com/daddydemir/crypto/config"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

var client *redis.Client

func GetRedisClient() *redis.Client {
	if client == nil {
		host := config.Get("REDIS_HOST")
		pass := config.Get("REDIS_PASS")

		client = redis.NewClient(&redis.Options{
			Addr:     host,
			Password: pass,
			DB:       0,
		})
	}

	if err := client.Ping(context.Background()).Err(); err != nil {
		slog.Error("Redis connection error: ", "error", err)
		panic(err)
	}

	return client
}
