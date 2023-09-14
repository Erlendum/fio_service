package cache

import (
	"context"
	"encoding/json"
	"fio_service/internal/config"
	"fio_service/pkg/cache"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(ctx context.Context, cfg config.RedisConfig) (cache.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (m *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := m.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return err
	}

	return nil
}

func (m *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	data, err := m.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var value interface{}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (m *RedisCache) Delete(ctx context.Context, key ...string) error {
	err := m.client.Del(ctx, key...).Err()
	if err != nil {
		return err
	}

	return nil
}
