package cacheadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
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

func (r *RedisCache) GetList(key string, list any, start, end int64) error {

	v := reflect.ValueOf(list)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("list must be a pointer to a slice")
	}

	slices := v.Elem()
	elemType := slices.Type().Elem()

	result, err := r.client.LRange(context.Background(), key, start, end).Result()
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

func (r *RedisCache) DeleteListItem(key string, start, end int64) error {
	ctx := context.Background()

	listLength, err := r.client.LLen(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get list length: %w", err)
	}

	if start < 0 || end >= listLength || start > end {
		return fmt.Errorf("invalid range: start (%d), end (%d), list length (%d)", start, end, listLength)
	}

	if start > 0 {
		err = r.client.LTrim(ctx, key, 0, start-1).Err()
		if err != nil {
			return fmt.Errorf("failed to trim the beginning of the list: %w", err)
		}
	}

	if end < listLength-1 {
		err = r.client.LTrim(ctx, key, end+1, listLength-1).Err()
		if err != nil {
			return fmt.Errorf("failed to trim the end of the list: %w", err)
		}
	}

	return nil
}

func (r *RedisCache) DeleteLastItem(key string) error {
	ctx := context.Background()

	val, err := r.client.RPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("list is empty or does not exist")
		}
		return fmt.Errorf("failed to remove last item: %w", err)
	}

	fmt.Printf("Removed last item: %s\n", val)
	return nil
}
