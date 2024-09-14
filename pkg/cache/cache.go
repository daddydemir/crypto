package cache

import (
	"time"
)

type Cache interface {
	Set(key string, value any) error
	SetWithExpiration(key string, value any, exp time.Duration) error
	SetList(key string, list any, exp time.Duration) error
	Get(key string) (any, error)
	GetList(key string, list any) error
	Delete(key string) error
}
