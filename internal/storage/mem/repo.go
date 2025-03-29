package mem

import (
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

// MemRepo просто обёртка над InMemoryRepo (можно вернуть напрямую — дело вкуса)
type Repo struct {
	*InMemoryRepo
}

var _ repository.URLRepository = (*Repo)(nil)

func NewRepo(log logger.Logger) *Repo {
	return &Repo{
		InMemoryRepo: New(log),
	}
}
