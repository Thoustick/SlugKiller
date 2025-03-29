package server

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/cache"
	"github.com/Thoustick/SlugKiller/internal/handler"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/internal/service"
	"github.com/Thoustick/SlugKiller/pkg/logger"
	"github.com/gin-gonic/gin"
)

type (
	StorageFactory func(ctx context.Context, cfg *config.Config, log logger.Logger) (repository.URLRepository, error)
	CacheProvider  func(cfg *config.Config, log logger.Logger) (cache.URLCache, error)
)

type App struct {
	engine *gin.Engine
	cfg    *config.Config
}

func NewApp(
	cfg *config.Config,
	storageFactory StorageFactory,
	cacheProvider CacheProvider,
) (*App, error) {
	ctx := setupGracefulShutdown()
	log := setupLogger(cfg)

	cacheLayer, err := cacheProvider(cfg, log)
	if err != nil {
		log.Error("Failed to initialize cache", err, nil)
		return nil, err
	}

	repo, err := storageFactory(ctx, cfg, log)
	if err != nil {
		log.Error("Failed to init sorage", err, nil)
		return nil, err
	}

	slugGen := service.NewSlugGenerator(cfg.SlugLength)
	urlService := service.NewURLService(repo, log, cacheLayer, cfg, slugGen)

	h := handler.NewHandler(urlService, log)

	r := setupRouter(h)

	return &App{
		engine: r,
		cfg:    cfg,
	}, nil
}

func (a *App) Run() error {
	return a.engine.Run(a.cfg.HTTPAddr)
}

func setupGracefulShutdown() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	_ = stop // использовать stop() можно в будущем, если нужно graceful stop
	return ctx
}

func setupLogger(cfg *config.Config) logger.Logger {
	return logger.InitLogger(cfg)
}

func setupRouter(h handler.URLHandler) *gin.Engine {
	r := gin.Default()
	h.RegisterRoutes(r)
	return r
}
