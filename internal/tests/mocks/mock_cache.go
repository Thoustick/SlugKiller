package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(ctx context.Context, slug string) (string, error) {
	args := m.Called(ctx, slug)
	return args.String(0), args.Error(1)
}

func (m *MockCache) SetNX(ctx context.Context, slug, url string, ttl time.Duration) error {
	args := m.Called(ctx, slug, url, ttl)
	return args.Error(0)
}
