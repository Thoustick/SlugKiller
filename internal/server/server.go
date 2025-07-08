package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/cache"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

type (
	StorageFactory func(ctx context.Context, cfg *config.Config, log logger.Logger) (repository.URLRepository, error)
	CacheProvider  func(cfg *config.Config, log logger.Logger) (cache.URLCache, error)
)

type App struct {
	Engine *gin.Engine
	Cfg    *config.Config
	Ctx    context.Context
	Cancel context.CancelFunc
	Logger logger.Logger
}

func (a *App) Run() error {
	srv := &http.Server{
		Addr:    a.Cfg.HTTPAddr,
		Handler: a.Engine,
	}

	go func() {
		a.Logger.Info("starting HTTP server", map[string]interface{}{
			"addr": a.Cfg.HTTPAddr,
		})

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Error("HTTP server error", err, nil)
		}
	}()

	<-a.Ctx.Done()

	a.Logger.Info("shutdown signal received", nil)

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		a.Logger.Error("graceful shutdown failed", err, nil)
		return err
	}

	a.Logger.Info("server stopped gracefully", nil)
	return nil
}
