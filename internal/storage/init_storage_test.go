package storage_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/storage"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

func TestInitStorage_Memory(t *testing.T) {
	log := logger.InitLogger(&config.Config{LogLevel: "debug"})
	cfg := &config.Config{
		StorageType: "memory",
	}

	repo, err := storage.InitStorage(context.Background(), cfg, log)
	assert.NoError(t, err)
	assert.NotNil(t, repo, "memory repo should not be nil")
}

func TestInitStorage_InvalidType(t *testing.T) {
	log := logger.InitLogger(&config.Config{LogLevel: "debug"})
	cfg := &config.Config{
		StorageType: "invalid",
	}

	repo, err := storage.InitStorage(context.Background(), cfg, log)
	assert.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "invalid storage type")
}

func TestInitStorage_Postgres(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set, skipping postgres storage test")
	}

	log := logger.InitLogger(&config.Config{LogLevel: "debug"})
	cfg := &config.Config{
		StorageType: "postgres",
		DatabaseURL: dsn,
		DBTimeout:   15,
	}

	repo, err := storage.InitStorage(context.Background(), cfg, log)
	assert.NoError(t, err)
	assert.NotNil(t, repo, "postgres repo should not be nil")
}
