package pg

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Thoustick/SlugKiller/internal/tests/mocks"
)

func TestGetBySlug(t *testing.T) {
	t.Run("успешное получение записи по slug", func(t *testing.T) {
		// 1. Создаём мок DBExecutor, мок Row, мок Logger
		dbMock := &mocks.MockDBExecutor{}
		rowMock := &mocks.MockRow{}
		loggerMock := &mocks.MockLogger{}

		// 2. Программируем rowMock, чтобы Scan заполнил поля model.Link
		rowMock.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			// args.Get(0) — это срез dest ...any
			dest := args.Get(0).([]any)
			*(dest[0].(*int64)) = 42                                              // link.ID
			*(dest[1].(*string)) = "test-slug"                                    // link.Slug
			*(dest[2].(*string)) = "https://example.com"                          // link.URL
			*(dest[3].(*time.Time)) = time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC) // link.CreatedAt
		}).Return(nil) // нет ошибки

		// 3. Программируем DBExecutor: QueryRow(...) вернёт rowMock
		dbMock.On("QueryRow", mock.Anything,
			"SELECT id, slug, url, created_at FROM urls WHERE slug = $1",
			[]interface{}{"test-slug"}).Return(rowMock)

		// loggerMock может вызываться или не вызываться. Здесь
		// нет ошибки, так что вызова Error(...) не ожидаем, но если нужно —
		// можно добавить проверку, что не вызывалось:
		loggerMock.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe()

		// 4. Собираем PostgresReader
		r := &PostgresReader{
			db:     dbMock,
			logger: loggerMock,
		}

		// 5. Тестируем
		link, err := r.GetBySlug(context.Background(), "test-slug")
		assert.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, int64(42), link.ID)
		assert.Equal(t, "test-slug", link.Slug)
		assert.Equal(t, "https://example.com", link.URL)
		assert.Equal(t, time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC), link.CreatedAt)

		// 6. Проверяем, что все ожидаемые вызовы произошли
		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})

	t.Run("ошибка: запись не найдена", func(t *testing.T) {
		// Допустим, если нет записи, Scan вернёт sql.ErrNoRows
		dbMock := &mocks.MockDBExecutor{}
		rowMock := &mocks.MockRow{}
		loggerMock := &mocks.MockLogger{}

		rowMock.On("Scan", mock.Anything).Return(errors.New("no rows in result set"))

		dbMock.On("QueryRow", mock.Anything,
			"SELECT id, slug, url, created_at FROM urls WHERE slug = $1",
			[]interface{}{"unknown"}).Return(rowMock)

		// Ожидаем, что logger.Error(...) будет вызван
		loggerMock.On("Error", "failed to get link by slug", mock.Anything, mock.Anything).Once()

		r := &PostgresReader{db: dbMock, logger: loggerMock}

		link, err := r.GetBySlug(context.Background(), "unknown")
		assert.Error(t, err)
		assert.Nil(t, link)

		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})

	t.Run("граничный случай: пустой slug", func(t *testing.T) {
		// Предположим, что пустой slug в БД не хранится,
		// и мы получим ошибку при сканировании.
		dbMock := &mocks.MockDBExecutor{}
		rowMock := &mocks.MockRow{}
		loggerMock := &mocks.MockLogger{}

		rowMock.On("Scan", mock.Anything).Return(errors.New("no rows in result set"))

		dbMock.On("QueryRow", mock.Anything,
			"SELECT id, slug, url, created_at FROM urls WHERE slug = $1",
			[]interface{}{""}).Return(rowMock)

		loggerMock.On("Error", "failed to get link by slug", mock.Anything, mock.Anything).Once()

		r := &PostgresReader{db: dbMock, logger: loggerMock}

		link, err := r.GetBySlug(context.Background(), "")
		assert.Error(t, err)
		assert.Nil(t, link)

		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})
}

func TestGetByOriginalURL(t *testing.T) {
	t.Run("успешное получение записи по URL", func(t *testing.T) {
		dbMock := &mocks.MockDBExecutor{}
		rowMock := &mocks.MockRow{}
		loggerMock := &mocks.MockLogger{}

		rowMock.On("Scan", mock.Anything).Run(func(args mock.Arguments) {
			dest := args.Get(0).([]any)
			*(dest[0].(*int64)) = 99
			*(dest[1].(*string)) = "slug-99"
			*(dest[2].(*string)) = "https://test.com"
			*(dest[3].(*time.Time)) = time.Date(2025, 3, 28, 0, 0, 0, 0, time.UTC)
		}).Return(nil)

		dbMock.On("QueryRow", mock.Anything,
			"SELECT id, slug, url, created_at FROM urls WHERE url = $1",
			[]interface{}{"https://test.com"}).Return(rowMock)

		r := &PostgresReader{db: dbMock, logger: loggerMock}

		link, err := r.GetByOriginalURL(context.Background(), "https://test.com")
		assert.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, int64(99), link.ID)
		assert.Equal(t, "slug-99", link.Slug)
		assert.Equal(t, "https://test.com", link.URL)
		assert.Equal(t, time.Date(2025, 3, 28, 0, 0, 0, 0, time.UTC), link.CreatedAt)

		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})

	t.Run("ошибка: запись не найдена", func(t *testing.T) {
		dbMock := &mocks.MockDBExecutor{}
		rowMock := &mocks.MockRow{}
		loggerMock := &mocks.MockLogger{}

		rowMock.On("Scan", mock.Anything).Return(errors.New("no rows"))
		dbMock.On("QueryRow", mock.Anything,
			"SELECT id, slug, url, created_at FROM urls WHERE url = $1",
			[]interface{}{"https://nope.com"}).Return(rowMock)

		loggerMock.On("Error", "failed to get link by original URL", mock.Anything, mock.Anything).Once()

		r := &PostgresReader{db: dbMock, logger: loggerMock}

		link, err := r.GetByOriginalURL(context.Background(), "https://nope.com")
		assert.Error(t, err)
		assert.Nil(t, link)

		dbMock.AssertExpectations(t)
		rowMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})
}
