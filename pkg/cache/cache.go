package cache

import (
	"github.com/daddydemir/crypto/pkg/cache/redis"
	"time"
)

type Cache interface {
	Set(key string, value any) error
	SetWithExpiration(key string, value any, exp time.Duration) error
	SetList(key string, list any, exp time.Duration) error
	Get(key string) (any, error)
	GetList(key string, list any, start, end int64) error
	DeleteListItem(key string, start, end int64) error
	DeleteLastItem(key string) error
	Delete(key string) error
}

func GetCacheService() Cache {
	return redis.NewRedisCache()
}
