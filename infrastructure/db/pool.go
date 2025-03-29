package db

import (
	"context"
	"fmt"
	"time"

	"github.com/Thoustick/SlugKiller/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PG struct {
	Pool   *pgxpool.Pool
	Logger logger.Logger
}

// New создаёт новое подключение к PostgreSQL
func New(ctx context.Context, connString string, timeOut time.Duration, log logger.Logger) (*PG, error) {
	ctxWTO, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	pool, err := pgxpool.New(ctxWTO, connString)
	if err != nil {
		log.Error("failed to connect to database", err, map[string]interface{}{
			"dsn": connString,
		})
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	log.Info("connected to PostgreSQL", nil)
	return &PG{Pool: pool, Logger: log}, nil
}

func (pg *PG) Close() {
	pg.Logger.Info("closing database connection", nil)
	pg.Pool.Close()
}
