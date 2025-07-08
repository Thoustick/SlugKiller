package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	Client *redis.Client
}

func NewRedisCache(client *redis.Client) URLCache {
	return &redisCache{Client: client}
}

func (r *redisCache) Get(ctx context.Context, slug string) (string, error) {
	val, err := r.Client.Get(ctx, slug).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *redisCache) SetNX(ctx context.Context, slug, url string, ttl time.Duration) error {
	return r.Client.SetNX(ctx, slug, url, ttl).Err()
}
