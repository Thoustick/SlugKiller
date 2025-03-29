package pg

import (
	"context"
	"errors"

	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/pkg/logger"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresWriter struct {
	db     DBExecutor
	logger logger.Logger
}

func NewPostgresWriter(db DBExecutor, l logger.Logger) *PostgresWriter {
	return &PostgresWriter{
		db:     db,
		logger: l,
	}
}

var _ repository.URLWriter = (*PostgresWriter)(nil)

func (w *PostgresWriter) Create(ctx context.Context, link *model.Link) error {
	const query = `INSERT INTO urls (slug, url, created_at) VALUES ($1, $2, $3)`
	_, err := w.db.Exec(ctx, query, link.Slug, link.URL, link.CreatedAt)
	if err != nil {

		// Обработка уникального конфликта (slug или url)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repository.ErrAlreadyExists
		}

		w.logger.Error("failed to insert link", err, map[string]interface{}{
			"slug": link.Slug,
			"url":  link.URL,
		})
		return err
	}
	return nil
}
