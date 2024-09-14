package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/daddydemir/crypto/config/cache"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"reflect"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		client: cache.GetRedisClient(),
	}
}

func (r *RedisCache) Set(key string, value any) error {
	response := r.client.Set(context.Background(), key, value, 0)
	return response.Err()
}

func (r *RedisCache) SetWithExpiration(key string, value any, expire time.Duration) error {
	response := r.client.Set(context.Background(), key, value, expire)
	return response.Err()
}

func (r *RedisCache) SetList(key string, list any, expire time.Duration) error {

	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Slice {
		return errors.New("list must be a slice")
	}

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()

		bytes, err := json.Marshal(elem)
		if err != nil {
			slog.Error("SetList:json.Marshal", "error", err)
			return err
		}

		err = r.client.RPush(context.Background(), key, bytes).Err()
		if err != nil {
			slog.Error("SetList:client.RPush", "error", err)
			return err
		}

		if expire > 0 {
			err = r.client.Expire(context.Background(), key, expire).Err()
			if err != nil {
				slog.Error("SetList:client.Expire", "error", err)
				return err
			}
		}
	}
	return nil
}

func (r *RedisCache) Get(key string) (any, error) {
	response := r.client.Get(context.Background(), key)
	result, _ := response.Result()
	return result, response.Err()
}

func (r *RedisCache) GetList(key string, list any) error {

	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("list must be a pointer to a slice")
	}

	slices := v.Elem()
	elemType := slices.Type().Elem()

	result, err := r.client.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		slog.Error("GetList:client.LRange", "error", err)
		return err
	}

	for _, jsonData := range result {
		newElem := reflect.New(elemType).Elem()
		err = json.Unmarshal([]byte(jsonData), newElem.Addr().Interface())
		if err != nil {
			slog.Error("GetList:json.Unmarshal", "error", err)
			return err
		}
		slices.Set(reflect.Append(slices, newElem))
	}

	return nil
}

func (r *RedisCache) Delete(key string) error {
	response := r.client.Del(context.Background(), key)
	return response.Err()
}
