package mocks

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
)

// MockDBExecutor - мок, реализующий интерфейс DBExecutor.
type MockDBExecutor struct {
	mock.Mock
}

func (m *MockDBExecutor) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	// В testify/mock обычно вызывается m.Called(...) для фиксации входных параметров
	// и возвращается значение, которое мы "программируем" в тесте.
	ret := m.Called(ctx, sql, args)
	// pgx.Row — это интерфейс, вернём заранее заготовленный объект (например, *MockRow)
	row, _ := ret.Get(0).(pgx.Row)
	return row
}

func (m *MockDBExecutor) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	ret := m.Called(ctx, sql, args)
	cmdTag, _ := ret.Get(0).(pgconn.CommandTag)
	err, _ := ret.Get(1).(error)
	return cmdTag, err
}
