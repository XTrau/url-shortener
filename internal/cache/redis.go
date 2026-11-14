package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

func NewRedisCache(rdb *redis.Client, timeout time.Duration) *RedisCache {
	return &RedisCache{
		rdb: rdb,
	}
}

func (kvs *RedisCache) Get(ctx context.Context, key string) (string, error) {
	res, err := kvs.rdb.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", ErrCacheKeyNotFound
	}

	return res, err
}

func (kvs *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return kvs.rdb.Set(ctx, key, value, ttl).Err()
}
