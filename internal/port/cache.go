package port

import "time"

type Cache interface {
	Set(key string, value any) error
	SetWithExpiration(key string, value any, expire time.Duration) error
	SetList(key string, list any, expire time.Duration) error

	Get(key string) (any, error)
	GetList(key string, list any, start, end int64) error

	Delete(key string) error
	DeleteListItem(key string, start, end int64) error
	DeleteLastItem(key string) error
}
