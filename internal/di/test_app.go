package di

import (
	"context"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/handler"
	"github.com/Thoustick/SlugKiller/internal/server"
	"github.com/Thoustick/SlugKiller/internal/service"
	"github.com/Thoustick/SlugKiller/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewTestApp(
	ctx context.Context,
	cfg *config.Config,
	storageFactory server.StorageFactory,
	cacheProvider server.CacheProvider,
) (*server.App, error) {
	log := logger.InitLogger(cfg)

	repo, err := storageFactory(ctx, cfg, log)
	if err != nil {
		return nil, err
	}

	cacheLayer, err := cacheProvider(cfg, log)
	if err != nil {
		return nil, err
	}

	slugGen := service.NewSlugGenerator(cfg.SlugLength)
	urlService := service.NewURLService(repo, log, cacheLayer, cfg, slugGen)

	h := handler.NewHandler(urlService, log)
	engine := gin.Default()
	h.RegisterRoutes(engine)

	return &server.App{
		Engine: engine,
		Cfg:    cfg,
		Ctx:    ctx,
		Logger: log,
	}, nil
}
