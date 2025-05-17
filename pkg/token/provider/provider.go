package provider

import (
	"context"
	"github.com/daddydemir/crypto/config/cache"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenProvider interface {
	GetValidToken() (string, error)
	MarkTokenAsExpired(token string, ttl time.Duration) error
}

type RedisTokenProvider struct {
	client     *redis.Client
	allKey     string
	expiredKey string
	ctx        context.Context
}

func NewRedisTokenProvider() *RedisTokenProvider {
	return &RedisTokenProvider{
		client:     cache.GetRedisClient(),
		allKey:     "all_tokens",
		expiredKey: "expired_tokens",
		ctx:        context.Background(),
	}
}

func (r *RedisTokenProvider) GetValidToken() (string, error) {
	allTokens, err := r.client.SMembers(r.ctx, r.allKey).Result()
	if err != nil {
		return "", err
	}

	expiredTokens, err := r.client.SMembers(r.ctx, r.expiredKey).Result()
	if err != nil {
		return "", err
	}
	expiredSet := make(map[string]struct{}, len(expiredTokens))
	for _, t := range expiredTokens {
		expiredSet[t] = struct{}{}
	}

	var validTokens []string
	for _, token := range allTokens {
		if _, expired := expiredSet[token]; !expired {
			validTokens = append(validTokens, token)
		}
	}

	if len(validTokens) == 0 {
		return "", nil
	}

	return validTokens[rand.Intn(len(validTokens))], nil
}

func (r *RedisTokenProvider) MarkTokenAsExpired(token string, ttl time.Duration) error {
	return r.client.SAdd(r.ctx, r.expiredKey, token).Err() // SAdd otomatik tekrar eklemeyi engeller
}
