package cache

import (
	"context"
	"time"
)

type URLCache interface {
	Get(ctx context.Context, slug string) (string, error)
	Set(ctx context.Context, slug, url string, ttl time.Duration) error
}
