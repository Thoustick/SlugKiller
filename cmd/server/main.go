package main

import (
	"context"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/di"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := di.NewApp(
		ctx,
		cfg,
	)
	if err != nil {
		panic(err)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
