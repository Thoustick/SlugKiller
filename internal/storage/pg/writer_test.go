// file: internal/storage/pg/writer_test.go
package pg

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Thoustick/SlugKiller/internal/model"
	"github.com/Thoustick/SlugKiller/internal/repository"
	"github.com/Thoustick/SlugKiller/internal/tests/mocks"
)

func TestCreate(t *testing.T) {
	t.Run("успешное создание ссылки", func(t *testing.T) {
		dbMock := &mocks.MockDBExecutor{}
		loggerMock := &mocks.MockLogger{}

		// Аргументы, которые мы будем передавать в Exec
		slug := "my-slug"
		url := "https://example.com"
		createdAt := time.Now()

		// Программируем Exec так, чтобы он возвращал успех (нет ошибки).
		// При успехе обычно возвращается какой-то CommandTag, например "INSERT 1".
		dbMock.On("Exec", mock.Anything,
			"INSERT INTO urls (slug, url, created_at) VALUES ($1, $2, $3)",
			[]interface{}{slug, url, createdAt},
		).Return(pgconn.NewCommandTag("INSERT 1"), nil).Once()

		// Мы НЕ ожидаем вызова loggerMock.Error(...) в случае успеха.
		// Если хотите, можете добавить:
		loggerMock.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe()

		writer := &PostgresWriter{
			db:     dbMock,
			logger: loggerMock,
		}
		link := &model.Link{
			Slug:      slug,
			URL:       url,
			CreatedAt: createdAt,
		}

		err := writer.Create(context.Background(), link)
		assert.NoError(t, err)

		dbMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})

	t.Run("ошибка: уникальный конфликт (ErrAlreadyExists)", func(t *testing.T) {
		dbMock := &mocks.MockDBExecutor{}
		loggerMock := &mocks.MockLogger{}

		slug := "duplicate-slug"
		url := "https://duplicate.com"
		createdAt := time.Now()

		// Чтобы смоделировать ошибку уникальности, вернём pgconn.PgError с кодом 23505
		pgErr := &pgconn.PgError{
			Code: "23505", // уникальный конфликт
		}
		dbMock.On("Exec", mock.Anything,
			"INSERT INTO urls (slug, url, created_at) VALUES ($1, $2, $3)",
			[]interface{}{slug, url, createdAt},
		).Return(pgconn.NewCommandTag(""), pgErr).Once()

		// В случае "23505" метод Create должен вернуть repository.ErrAlreadyExists
		// и НЕ должен логировать ошибку, если мы так решили (или может логировать – зависит от логики).
		// Допустим, код у вас явно возвращает `repository.ErrAlreadyExists`, минуя логгер.
		// В таком случае не ждём вызова loggerMock.Error(...).
		loggerMock.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe()

		writer := &PostgresWriter{db: dbMock, logger: loggerMock}
		link := &model.Link{
			Slug:      slug,
			URL:       url,
			CreatedAt: createdAt,
		}

		err := writer.Create(context.Background(), link)
		assert.ErrorIs(t, err, repository.ErrAlreadyExists)

		dbMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})

	t.Run("ошибка: другая ошибка БД", func(t *testing.T) {
		dbMock := &mocks.MockDBExecutor{}
		loggerMock := &mocks.MockLogger{}

		slug := "some-slug"
		url := "https://test.com"
		createdAt := time.Now()

		dbMock.On("Exec", mock.Anything,
			"INSERT INTO urls (slug, url, created_at) VALUES ($1, $2, $3)",
			[]interface{}{slug, url, createdAt},
		).Return(pgconn.NewCommandTag(""), errors.New("db failure")).Once()

		// В таком случае код должен вызвать logger.Error(...)
		loggerMock.On("Error", "failed to insert link", mock.Anything, mock.Anything).Once()

		writer := &PostgresWriter{db: dbMock, logger: loggerMock}
		link := &model.Link{Slug: slug, URL: url, CreatedAt: createdAt}

		err := writer.Create(context.Background(), link)
		// Ожидаем вернуть исходную ошибку (не ErrAlreadyExists)
		assert.Error(t, err)
		assert.NotEqual(t, repository.ErrAlreadyExists, err)

		dbMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})

	t.Run("граничный случай: пустое поле (slug или url)", func(t *testing.T) {
		// Зависит от того, как вы хотите это обрабатывать.
		// Предположим, что в БД стоит NOT NULL / UNIQUE constraint
		// и при попытке вставить пустые поля получим ошибку "23502" (NOT NULL violation),
		// или любую другую. Вариации зависят от структуры схемы.
		dbMock := &mocks.MockDBExecutor{}
		loggerMock := &mocks.MockLogger{}

		// Допустим, slug = "" => нарушение NOT NULL
		pgErr := &pgconn.PgError{
			Code:    "23502",
			Message: "null value in column \"slug\" of relation \"urls\" violates not-null constraint",
		}
		dbMock.On("Exec",
			mock.Anything,
			"INSERT INTO urls (slug, url, created_at) VALUES ($1, $2, $3)",
			mock.MatchedBy(func(args []interface{}) bool {
				return len(args) == 3 &&
					args[0] == "" && // slug
					args[1] == "https://gaps.com" // url
				// args[2] — неважно, пропускаем
			}),
		).Return(pgconn.NewCommandTag(""), pgErr).Once()

		loggerMock.On("Error", "failed to insert link", pgErr, mock.Anything).Once()

		writer := &PostgresWriter{db: dbMock, logger: loggerMock}
		link := &model.Link{
			Slug: "",
			URL:  "https://gaps.com",
			// CreatedAt: time.Now(),
		}

		err := writer.Create(context.Background(), link)
		assert.Error(t, err)
		// Это не ошибка "23505", так что ожидаем, что вернётся именно pgErr
		// (или обёрнутое что-то), но код всё равно вызывал logger.Error(...)
		dbMock.AssertExpectations(t)
		loggerMock.AssertExpectations(t)
	})
}
