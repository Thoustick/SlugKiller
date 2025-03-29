package redis

import (
	"context"
	"fmt"

	"github.com/Thoustick/SlugKiller/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
	logger logger.Logger
}

// NewRedisClient инициализирует Redis клиента
func NewRedisClient(redisURL, password string, dbIndex int, l logger.Logger) (*Client, error) {
	ctx := context.Background()
	l.Info("Initializing Redis connection", map[string]interface{}{
		"host": redisURL,
		"db":   dbIndex,
	})

	client := redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: password,
		DB:       dbIndex,
	})

	if _, err := client.Ping(ctx).Result(); err != nil {
		l.Error("Unable to connect to Redis", err, map[string]interface{}{
			"host": redisURL,
			"db":   dbIndex,
		})
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	l.Info("Successfully connected to Redis", map[string]interface{}{
		"host": redisURL,
		"db":   dbIndex,
	})

	return &Client{
		client: client,
		ctx:    ctx,
		logger: l,
	}, nil
}

func (r *Client) Client() *redis.Client {
	return r.client
}
