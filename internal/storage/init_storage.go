package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/infrastructure/db"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/internal/storage/mem"
	"github.com/Thoustick/SlugKiller/internal/storage/pg"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

func InitStorage(ctx context.Context, cfg *config.Config, log logger.Logger) (repository.URLRepository, error) {
	switch cfg.StorageType {
	case "postgres":
		if cfg.DatabaseURL == "" {
			return nil, errors.New("DATABASE_URL not set")
		}
		db, err := db.New(ctx, cfg.DatabaseURL, cfg.DBTimeout, log)
		if err != nil {
			return nil, err
		}
		return pg.NewRepo(db.Pool, log), nil

	case "memory":
		return mem.NewRepo(log), nil
	default:
		return nil, fmt.Errorf("invalid storage type: %s", cfg.StorageType)
	}
}
