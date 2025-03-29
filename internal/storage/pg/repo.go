package pg

import (
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

// PostgresRepo объединяет ридер и райтер в один объект
type PostgresRepo struct {
	*PostgresReader
	*PostgresWriter
}

// Проверка реализации интерфейса
var _ repository.URLRepository = (*PostgresRepo)(nil)

// NewRepo создаёт единый PostgresRepo
func NewRepo(db DBExecutor, log logger.Logger) *PostgresRepo {
	return &PostgresRepo{
		PostgresReader: NewPostgresReader(db, log),
		PostgresWriter: NewPostgresWriter(db, log),
	}
}
