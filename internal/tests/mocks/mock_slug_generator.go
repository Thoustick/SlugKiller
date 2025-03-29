package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockSlugGenerator struct {
	mock.Mock
}

func (m *MockSlugGenerator) Generate(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}
