package mocks

import (
	"context"

	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/stretchr/testify/mock"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) GetByOriginalURL(ctx context.Context, original string) (*model.Link, error) {
	args := m.Called(ctx, original)
	link := args.Get(0)
	if link == nil {
		return nil, args.Error(1)
	}
	return link.(*model.Link), args.Error(1)
}

func (m *MockURLRepository) GetBySlug(ctx context.Context, slug string) (*model.Link, error) {
	args := m.Called(ctx, slug)
	link := args.Get(0)
	if link == nil {
		return nil, args.Error(1)
	}
	return link.(*model.Link), args.Error(1)
}

func (m *MockURLRepository) Create(ctx context.Context, url *model.Link) error {
	args := m.Called(ctx, url)
	return args.Error(0)
}

func (m *MockURLRepository) Delete(ctx context.Context, slug string) error {
	args := m.Called(ctx, slug)
	return args.Error(0)
}
