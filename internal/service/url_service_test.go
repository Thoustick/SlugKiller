package service_test

import (
	"context"
	"errors"
	"testing"
	"unicode"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/internal/service"
	"github.com/Thoustick/SlugKiller/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testURLService struct {
	svc     service.URLService
	repo    *mocks.MockURLRepository
	cache   *mocks.MockCache
	logger  *mocks.MockLogger
	slugGen *mocks.MockSlugGenerator
}

func setupURLService() testURLService {
	repo := new(mocks.MockURLRepository)
	cache := new(mocks.MockCache)
	logger := new(mocks.MockLogger)
	slugGen := new(mocks.MockSlugGenerator)

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe()
	logger.On("Warn", mock.Anything, mock.Anything).Maybe()
	logger.On("Info", mock.Anything, mock.Anything).Maybe()
	logger.On("Debug", mock.Anything, mock.Anything).Maybe()
	logger.On("Fatal", mock.Anything, mock.Anything, mock.Anything).Maybe()

	cfg := &config.Config{CacheTTL: 0, MaxAttempts: 5, SlugLength: 10}
	svc := service.NewURLService(repo, logger, cache, cfg, slugGen)

	return testURLService{
		svc:     svc,
		repo:    repo,
		cache:   cache,
		logger:  logger,
		slugGen: slugGen,
	}
}

func TestShorten_EmptyURL(t *testing.T) {
	ts := setupURLService()
	slug, err := ts.svc.Shorten(context.Background(), "")

	assert.Error(t, err)
	assert.Empty(t, slug)
	ts.repo.AssertExpectations(t)
	ts.slugGen.AssertExpectations(t)
}

func TestShorten_ExistingURL(t *testing.T) {
	ts := setupURLService()
	original := "https://example.com"
	expectedSlug := "abc123"

	ts.repo.On("GetByOriginalURL", mock.Anything, original).Return(&model.Link{
		URL:  original,
		Slug: expectedSlug,
	}, nil)

	slug, err := ts.svc.Shorten(context.Background(), original)

	assert.NoError(t, err)
	assert.Equal(t, expectedSlug, slug)
	ts.repo.AssertExpectations(t)
}

func TestShorten_DBError(t *testing.T) {
	ts := setupURLService()
	original := "https://example.com"
	dbErr := errors.New("db down")

	ts.repo.On("GetByOriginalURL", mock.Anything, original).Return(nil, dbErr)

	slug, err := ts.svc.Shorten(context.Background(), original)

	assert.Error(t, err)
	assert.Empty(t, slug)
	ts.repo.AssertExpectations(t)
}

func TestShorten_CreateNewSlug(t *testing.T) {
	ts := setupURLService()
	original := "https://newsite.com"
	generatedSlug := "customSlug123"

	ts.repo.On("GetByOriginalURL", mock.Anything, original).Return(nil, repository.ErrNotFound)
	ts.slugGen.On("Generate", mock.Anything).Return(generatedSlug, nil).Once()
	ts.repo.On("GetBySlug", mock.Anything, generatedSlug).Return(nil, repository.ErrNotFound).Once()
	ts.repo.On("Create", mock.Anything, mock.AnythingOfType("*model.Link")).Return(nil).Once()

	slug, err := ts.svc.Shorten(context.Background(), original)

	assert.NoError(t, err)
	assert.Equal(t, generatedSlug, slug)
	ts.repo.AssertExpectations(t)
	ts.slugGen.AssertExpectations(t)
}

func TestShorten_CreateNewSlug_Retries(t *testing.T) {
	ts := setupURLService()
	original := "https://retrytest.com"

	firstSlug := "slug_collision"
	secondSlug := "slug_collision_again"
	finalSlug := "slug_unique"

	ts.repo.On("GetByOriginalURL", mock.Anything, original).Return(nil, repository.ErrNotFound)

	// Эмуляция трёх попыток с коллизиями
	ts.slugGen.On("Generate", mock.Anything).Return(firstSlug, nil).Once()
	ts.repo.On("GetBySlug", mock.Anything, firstSlug).Return(nil, nil).Once() // Занят
	ts.slugGen.On("Generate", mock.Anything).Return(secondSlug, nil).Once()
	ts.repo.On("GetBySlug", mock.Anything, secondSlug).Return(nil, nil).Once() // Тоже занят

	ts.slugGen.On("Generate", mock.Anything).Return(finalSlug, nil).Once()
	ts.repo.On("GetBySlug", mock.Anything, finalSlug).Return(nil, repository.ErrNotFound).Once()
	ts.repo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

	slug, err := ts.svc.Shorten(context.Background(), original)

	assert.NoError(t, err)
	assert.Equal(t, finalSlug, slug)
	ts.repo.AssertExpectations(t)
	ts.slugGen.AssertExpectations(t)
}

func TestShorten_GenerationFails(t *testing.T) {
	ts := setupURLService()
	original := "https://error-during-generation.com"
	genErr := errors.New("slug generation error")

	ts.repo.On("GetByOriginalURL", mock.Anything, original).Return(nil, repository.ErrNotFound)
	ts.slugGen.On("Generate", mock.Anything).Return("", genErr).Once()

	slug, err := ts.svc.Shorten(context.Background(), original)
	assert.Error(t, err)
	assert.Empty(t, slug)
	ts.repo.AssertExpectations(t)
	ts.slugGen.AssertExpectations(t)
}

func setupResolveService() (service.URLService, *mocks.MockURLRepository, *mocks.MockCache, *mocks.MockLogger) {
	repo := new(mocks.MockURLRepository)
	cache := new(mocks.MockCache)
	logger := new(mocks.MockLogger)
	slugGen := new(mocks.MockSlugGenerator)

	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe()
	logger.On("Warn", mock.Anything, mock.Anything).Maybe()
	logger.On("Info", mock.Anything, mock.Anything).Maybe()
	logger.On("Debug", mock.Anything, mock.Anything).Maybe()
	logger.On("Fatal", mock.Anything, mock.Anything, mock.Anything).Maybe()

	cfg := &config.Config{CacheTTL: 0}
	svc := service.NewURLService(repo, logger, cache, cfg, slugGen)
	return svc, repo, cache, logger
}

func TestResolve_EmptySlug(t *testing.T) {
	svc, _, _, _ := setupResolveService()
	url, err := svc.Resolve(context.Background(), "")
	assert.Error(t, err)
	assert.Empty(t, url)
}

func TestResolve_CacheHit(t *testing.T) {
	svc, _, cache, _ := setupResolveService()
	slug := "abc123"
	originalURL := "https://example.com"

	cache.On("Get", mock.Anything, slug).Return(originalURL, nil)

	url, err := svc.Resolve(context.Background(), slug)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)
	cache.AssertExpectations(t)
}

func TestResolve_CacheMiss_DBHit(t *testing.T) {
	svc, repo, cache, _ := setupResolveService()
	slug := "abc123"
	originalURL := "https://example.com"

	cache.On("Get", mock.Anything, slug).Return("", errors.New("cache miss"))
	repo.On("GetBySlug", mock.Anything, slug).Return(&model.Link{Slug: slug, URL: originalURL}, nil)
	cache.On("Set", mock.Anything, slug, originalURL, mock.Anything).Return(nil)

	url, err := svc.Resolve(context.Background(), slug)

	assert.NoError(t, err)
	assert.Equal(t, originalURL, url)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestResolve_SlugNotFound(t *testing.T) {
	svc, repo, cache, _ := setupResolveService()
	slug := "notfound"

	cache.On("Get", mock.Anything, slug).Return("", errors.New("cache miss"))
	repo.On("GetBySlug", mock.Anything, slug).Return(nil, repository.ErrNotFound)

	url, err := svc.Resolve(context.Background(), slug)

	assert.Error(t, err)
	assert.Empty(t, url)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestResolve_DBError(t *testing.T) {
	svc, repo, cache, _ := setupResolveService()
	slug := "oops"
	dbErr := errors.New("db failure")

	cache.On("Get", mock.Anything, slug).Return("", errors.New("cache miss"))
	repo.On("GetBySlug", mock.Anything, slug).Return(nil, dbErr)

	url, err := svc.Resolve(context.Background(), slug)

	assert.Error(t, err)
	assert.Empty(t, url)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestDefaultSlugGenerator_Generate(t *testing.T) {
	gen := service.NewSlugGenerator(10)
	slug, err := gen.Generate(context.Background())

	assert.NoError(t, err)
	assert.Len(t, slug, 10)

	for _, r := range slug {
		isValid := unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
		assert.True(t, isValid, "unexpected rune in slug: %q", r)
	}
}
