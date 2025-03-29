package service

import "context"

type URLService interface {
	Shorten(ctx context.Context, originalURL string) (string, error)
	Resolve(ctx context.Context, shortURL string) (string, error)
}

type SlugGenerator interface {
	Generate(ctx context.Context) (string, error)
}
