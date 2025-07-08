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

func NewApp(ctx context.Context, cfg *config.Config) (*server.App, error) {
	log := logger.InitLogger(cfg)

	// Инициализация хранилища
	repo, err := server.ProductionStorageFactory(ctx, cfg, log)
	if err != nil {
		log.Error("failed to initialize storage", err, nil)
		return nil, err
	}

	// Инициализация кэша
	cacheLayer, err := server.ProductionCacheProvider(cfg, log)
	if err != nil {
		log.Error("failed to initialize cache", err, nil)
		return nil, err
	}

	slugGen := service.NewSlugGenerator(cfg.SlugLength)

	urlServiceInstance := service.NewURLService(
		repo,
		log,
		cacheLayer,
		cfg,
		slugGen,
	)

	h := handler.NewHandler(urlServiceInstance, log)

	r := setupRouter(h)

	appCtx, cancel := context.WithCancel(ctx)

	return &server.App{
		Engine: r,
		Cfg:    cfg,
		Ctx:    appCtx,
		Cancel: cancel,
		Logger: log,
	}, nil
}

func setupRouter(h handler.URLHandler) *gin.Engine {
	r := gin.Default()  // <- Инициализация маршрутизатора Gin
	h.RegisterRoutes(r) // <- Регистрируем маршруты
	return r            // <- Возвращаем готовый engine
}
