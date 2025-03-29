package mocks

import (
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

// MockRow - мок, реализующий pgx.Row.
type MockRow struct {
	mock.Mock
}

var _ pgx.Row = (*MockRow)(nil)

func (r *MockRow) Scan(dest ...any) error {
	// Вызываем .Called и получаем что вернёт Testify.
	ret := r.Called(dest)
	return ret.Error(0)
}
