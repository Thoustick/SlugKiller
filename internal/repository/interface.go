package repository

import (
	"context"

	"github.com/Thoustick/SlugKiller/internal/model"
)

// URLReader defines read-only operations for URL entities.
type URLReader interface {
	GetBySlug(ctx context.Context, slug string) (*model.Link, error)
	GetByOriginalURL(ctx context.Context, original string) (*model.Link, error)
}

// URLWriter defines write operations for URL entities.
type URLWriter interface {
	Create(ctx context.Context, url *model.Link) error
}

// URLRepository combines read and write operations.
type URLRepository interface {
	URLReader
	URLWriter
}
