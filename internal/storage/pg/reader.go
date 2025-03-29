package pg

import (
	"context"

	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

type PostgresReader struct {
	db     DBExecutor
	logger logger.Logger
}

func NewPostgresReader(db DBExecutor, l logger.Logger) *PostgresReader {
	return &PostgresReader{
		db:     db,
		logger: l,
	}
}

var _ repository.URLReader = (*PostgresReader)(nil)

func (r *PostgresReader) GetBySlug(ctx context.Context, slug string) (*model.Link, error) {
	const query = `SELECT id, slug, url, created_at FROM urls WHERE slug = $1`
	var link model.Link

	err := r.db.QueryRow(ctx, query, slug).Scan(&link.ID, &link.Slug, &link.URL, &link.CreatedAt)
	if err != nil {
		r.logger.Error("failed to get link by slug", err, map[string]interface{}{
			"slug": slug,
		})
		return nil, err
	}

	return &link, nil
}

func (r *PostgresReader) GetByOriginalURL(ctx context.Context, url string) (*model.Link, error) {
	const query = `SELECT id, slug, url, created_at FROM urls WHERE url = $1`
	var link model.Link

	err := r.db.QueryRow(ctx, query, url).Scan(&link.ID, &link.Slug, &link.URL, &link.CreatedAt)
	if err != nil {
		r.logger.Error("failed to get link by original URL", err, map[string]interface{}{
			"url": url,
		})
		return nil, err
	}

	return &link, nil
}
