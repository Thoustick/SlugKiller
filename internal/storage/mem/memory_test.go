package mem_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/internal/storage/mem"
	"github.com/Thoustick/SlugKiller/internal/tests/mocks"
)

func TestInMemoryRepo_Create_And_Get(t *testing.T) {
	log := new(mocks.MockLogger)
	repo := mem.NewRepo(log)
	ctx := context.Background()

	link := &model.Link{
		Slug: "abc123",
		URL:  "https://example.com",
	}

	err := repo.Create(ctx, link)
	require.NoError(t, err)

	gotBySlug, err := repo.GetBySlug(ctx, "abc123")
	assert.NoError(t, err)
	assert.Equal(t, link.URL, gotBySlug.URL)

	gotByOriginal, err := repo.GetByOriginalURL(ctx, "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, link.Slug, gotByOriginal.Slug)
}

func TestInMemoryRepo_GetBySlug_NotFound(t *testing.T) {
	log := new(mocks.MockLogger)
	repo := mem.NewRepo(log)
	ctx := context.Background()

	_, err := repo.GetBySlug(ctx, "not_exist")
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestInMemoryRepo_GetByOriginalURL_NotFound(t *testing.T) {
	log := new(mocks.MockLogger)
	repo := mem.NewRepo(log)
	ctx := context.Background()

	_, err := repo.GetByOriginalURL(ctx, "https://notfound.com")
	assert.ErrorIs(t, err, repository.ErrNotFound)
}

func TestInMemoryRepo_Create_DuplicateSlug(t *testing.T) {
	log := new(mocks.MockLogger)
	repo := mem.NewRepo(log)
	ctx := context.Background()

	link := &model.Link{
		Slug: "abc123",
		URL:  "https://url1.com",
	}
	_ = repo.Create(ctx, link)

	dup := &model.Link{
		Slug: "abc123",
		URL:  "https://url2.com",
	}
	err := repo.Create(ctx, dup)
	assert.ErrorIs(t, err, repository.ErrAlreadyExists)
}

func TestInMemoryRepo_Create_DuplicateURL(t *testing.T) {
	log := new(mocks.MockLogger)
	repo := mem.NewRepo(log)
	ctx := context.Background()

	link := &model.Link{
		Slug: "slug1",
		URL:  "https://dupe.com",
	}
	_ = repo.Create(ctx, link)

	dup := &model.Link{
		Slug: "slug2",
		URL:  "https://dupe.com",
	}
	err := repo.Create(ctx, dup)
	assert.ErrorIs(t, err, repository.ErrAlreadyExists)
}
