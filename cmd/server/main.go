package main

import (
	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/server"
)

func main() {
	cfg := config.Load()
	app, err := server.NewApp(
		cfg,
		server.ProductionStorageFactory,
		server.ProductionCacheProvider,
	)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
