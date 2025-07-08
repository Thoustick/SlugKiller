package service

import (
	"context"
	"errors"
	"time"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/cache"
	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/pkg/logger"
)

type urlService struct {
	repo    repository.URLRepository
	cache   cache.URLCache
	logger  logger.Logger
	cfg     *config.Config
	slugGen SlugGenerator
}

func NewURLService(
	r repository.URLRepository,
	l logger.Logger,
	c cache.URLCache,
	cfg *config.Config,
	slugGen SlugGenerator,
) URLService {
	return &urlService{
		repo:    r,
		cache:   c,
		logger:  l,
		cfg:     cfg,
		slugGen: slugGen,
	}
}

func (s *urlService) CreateUniqueSlugLoop(ctx context.Context, originalURL string) (string, error) {
	for i := 0; i < s.cfg.MaxAttempts; i++ {
		slug, err := s.generateUniqueSlug(ctx)
		if err != nil {
			s.logger.Error("Failed to generate slug", err, map[string]interface{}{
				"attempt": i + 1,
				"url":     originalURL,
			})
			return "", err
		}

		link := &model.Link{
			Slug:      slug,
			URL:       originalURL,
			CreatedAt: time.Now(),
		}

		err = s.repo.Create(ctx, link)
		if err == nil {
			return slug, nil
		}
		if errors.Is(err, repository.ErrAlreadyExists) {
			s.logger.Warn("Slug already exists, retrying", map[string]interface{}{
				"slug": slug,
				"try":  i + 1,
			})
			continue
		}

		s.logger.Error("Failed to create link", err, map[string]interface{}{
			"url":  originalURL,
			"slug": slug,
		})
		return "", err
	}
	return "", errors.New("failed to generate unique slug after max attempts")
}

// Shorten проверяет, есть ли уже запись для originalURL.
// Если есть, возвращает существующий slug.
// Если нет, генерирует уникальный slug и сохраняет новую запись в базе.
func (s *urlService) Shorten(ctx context.Context, originalURL string) (string, error) {
	if originalURL == "" {
		s.logger.Warn("Attempted to shorten empty URL", nil)
		return "", errors.New("empty URL provided")
	}

	// 1. Проверяем, нет ли уже записи
	existingLink, err := s.repo.GetByOriginalURL(ctx, originalURL)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) { // <--- вот правильная проверка
			s.logger.Error("Failed to check existing URL", err, map[string]interface{}{
				"url": originalURL,
			})
			return "", err
		}
		existingLink = nil // просто URL ещё не сокращали, это нормально
	}

	// 2. Если запись уже есть — возвращаем slug
	if existingLink != nil {
		s.logger.Info("URL already shortened, returning existing slug", map[string]interface{}{
			"url":  originalURL,
			"slug": existingLink.Slug,
		})
		return existingLink.Slug, nil
	}

	// 3. Генерируем уникальный slug
	slug, err := s.CreateUniqueSlugLoop(ctx, originalURL)
	if err != nil {
		return "", err
	}

	s.logger.Info("Successfully shortened URL", map[string]interface{}{
		"url":  originalURL,
		"slug": slug,
	})
	return slug, nil
}

// service/url_service.go
func (s *urlService) Resolve(ctx context.Context, slug string) (string, error) {
	if slug == "" {
		s.logger.Warn("Empty slug in resolve", nil)
		return "", errors.New("empty slug")
	}

	// 1. Проверяем кэш
	url, err := s.cache.Get(ctx, slug)
	if err == nil {
		s.logger.Info("Cache hit", map[string]interface{}{"slug": slug})
		return url, nil
	}

	// Логируем только НЕ "cache miss" ошибки
	if !errors.Is(err, cache.ErrCacheMiss) {
		s.logger.Warn("Cache lookup failed", map[string]interface{}{
			"slug":  slug,
			"error": err.Error(),
		})
	}

	// 2. Идём в базу
	link, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			s.logger.Warn("Slug not found", map[string]interface{}{"slug": slug})
			return "", err
		}
		s.logger.Error("Failed to fetch slug from DB", err, nil)
		return "", err
	}

	// 3. Обновляем кэш (добавляем обработку ошибок записи)
	if err := s.cache.SetNX(ctx, slug, link.URL, s.cfg.CacheTTL); err != nil {
		s.logger.Warn("Failed to update cache", map[string]interface{}{
			"slug":  slug,
			"error": err.Error(),
		})
	}

	s.logger.Info("Slug resolved and cached", map[string]interface{}{
		"slug": slug,
		"url":  link.URL,
	})
	return link.URL, nil
}
