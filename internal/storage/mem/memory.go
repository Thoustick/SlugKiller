package mem

import (
	"context"
	"sync"
	"time"

	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

type InMemoryRepo struct {
	mu       sync.RWMutex
	bySlug   map[string]*model.Link
	byOrigin map[string]*model.Link
	logger   logger.Logger
}

// New возвращает in-memory хранилище, реализующее URLRepository
func New(log logger.Logger) *InMemoryRepo {
	return &InMemoryRepo{
		bySlug:   make(map[string]*model.Link),
		byOrigin: make(map[string]*model.Link),
		logger:   log,
	}
}

// Проверка реализации интерфейса
var _ repository.URLRepository = (*InMemoryRepo)(nil)

func (r *InMemoryRepo) GetBySlug(_ context.Context, slug string) (*model.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.bySlug[slug]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return link, nil
}

func (r *InMemoryRepo) GetByOriginalURL(_ context.Context, original string) (*model.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.byOrigin[original]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return link, nil
}

func (r *InMemoryRepo) Create(_ context.Context, link *model.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.bySlug[link.Slug]; exists {
		return repository.ErrAlreadyExists
	}
	if _, exists := r.byOrigin[link.URL]; exists {
		return repository.ErrAlreadyExists
	}

	link.ID = int64(len(r.bySlug) + 1)
	link.CreatedAt = time.Now()

	r.bySlug[link.Slug] = link
	r.byOrigin[link.URL] = link
	return nil
}
