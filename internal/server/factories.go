package server

import (
	"context"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/infrastructure/redis"
	"github.com/Thoustick/SlugKiller/internal/cache"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/internal/storage"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

// Реальная фабрика для продакшена
func ProductionStorageFactory(ctx context.Context, cfg *config.Config, log logger.Logger) (repository.URLRepository, error) {
	return storage.InitStorage(ctx, cfg, log)
}

func ProductionCacheProvider(cfg *config.Config, log logger.Logger) (cache.URLCache, error) {
	redisClient, err := redis.NewRedisClient(
		cfg.RedisHost,
		cfg.RedisPass,
		cfg.RedisDB,
		log,
	)
	if err != nil {
		return nil, err
	}
	return cache.NewRedisCache(redisClient.Client()), nil
}
