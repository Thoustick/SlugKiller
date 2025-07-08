package server_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/cache"
	"github.com/Thoustick/SlugKiller/internal/di"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/internal/tests/mocks"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

func TestNewAppWithMocks(t *testing.T) {
	t.Run("успешная инициализация", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{
			HTTPAddr:    ":8080",
			StorageType: "memory",
			SlugLength:  8,
			MaxAttempts: 5,
			CacheTTL:    300,
		}

		mockStorage := func(_ context.Context, _ *config.Config, _ logger.Logger) (repository.URLRepository, error) {
			return &mocks.MockURLRepository{}, nil
		}
		mockCache := func(_ *config.Config, _ logger.Logger) (cache.URLCache, error) {
			return &mocks.MockCache{}, nil
		}

		app, err := di.NewTestApp(ctx, cfg, mockStorage, mockCache)
		assert.NoError(t, err)
		assert.NotNil(t, app)
		assert.Equal(t, ":8080", cfg.HTTPAddr)
	})

	t.Run("ошибка при инициализации cache", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{
			HTTPAddr:    ":8080",
			StorageType: "memory",
			SlugLength:  8,
			MaxAttempts: 5,
			CacheTTL:    300,
		}

		mockStorage := func(_ context.Context, _ *config.Config, _ logger.Logger) (repository.URLRepository, error) {
			return &mocks.MockURLRepository{}, nil
		}
		mockCache := func(_ *config.Config, _ logger.Logger) (cache.URLCache, error) {
			return nil, errors.New("cache init fail")
		}

		app, err := di.NewTestApp(ctx, cfg, mockStorage, mockCache)
		assert.Nil(t, app)
		assert.EqualError(t, err, "cache init fail")
	})

	t.Run("ошибка при инициализации storage", func(t *testing.T) {
		ctx := context.Background()
		cfg := &config.Config{
			HTTPAddr:    ":8080",
			StorageType: "memory",
			SlugLength:  8,
			MaxAttempts: 5,
			CacheTTL:    300,
		}

		mockStorage := func(_ context.Context, _ *config.Config, _ logger.Logger) (repository.URLRepository, error) {
			return nil, errors.New("db down")
		}
		mockCache := func(_ *config.Config, _ logger.Logger) (cache.URLCache, error) {
			return &mocks.MockCache{}, nil
		}

		app, err := di.NewTestApp(ctx, cfg, mockStorage, mockCache)
		assert.Nil(t, app)
		assert.EqualError(t, err, "db down")
	})
}

func TestApp_Run(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		HTTPAddr:    ":8081",
		StorageType: "memory",
		SlugLength:  8,
		MaxAttempts: 5,
		CacheTTL:    60,
	}

	mockStorage := func(_ context.Context, _ *config.Config, _ logger.Logger) (repository.URLRepository, error) {
		return &mocks.MockURLRepository{}, nil
	}

	mockCache := func(_ *config.Config, _ logger.Logger) (cache.URLCache, error) {
		return &mocks.MockCache{}, nil
	}

	app, err := di.NewTestApp(ctx, cfg, mockStorage, mockCache)
	assert.NoError(t, err)

	// Запускаем сервер в отдельной горутине
	go func() {
		err := app.Run()
		assert.NoError(t, err)
	}()

	// Даем серверу время на старт
	time.Sleep(100 * time.Millisecond)

	// Пытаемся сделать POST-запрос на /shorten
	resp, err := http.Post("http://localhost:8081/shorten", "application/json", bytes.NewBuffer([]byte(`{}`)))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
