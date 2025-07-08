package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Thoustick/SlugKiller/internal/handler"
	"github.com/Thoustick/SlugKiller/internal/tests/mocks"
)

// request структурки (если нужны)
type ShortenRequest struct {
	URL string `json:"url"`
}

func TestShortenURL_Success(t *testing.T) {
	// 1. Мокаем зависимости
	svc := new(mocks.MockURLService)
	log := new(mocks.MockLogger)
	h := handler.NewHandler(svc, log)

	// 2. Настраиваем мок
	testURL := "https://example.com"
	testSlug := "abc123"
	svc.On("Shorten", mock.Anything, testURL).Return(testSlug, nil)

	// Логгер можем не проверять досконально
	log.On("Info", mock.Anything, mock.Anything).Maybe()
	log.On("Error", mock.Anything, mock.Anything, mock.Anything).Maybe()
	log.On("Warn", mock.Anything, mock.Anything).Maybe()

	// 3. Создаём тестовый рекордер и контекст Gin
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 4. Имитируем POST-запрос
	body := `{"url":"` + testURL + `"}`
	req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// 5. Вызываем хендлер
	h.ShortenURL(c)

	// 6. Проверяем ответ
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), testSlug) // {"slug":"abc123"}
	svc.AssertExpectations(t)
	log.AssertExpectations(t)
}

func TestShortenURL_InvalidJSON(t *testing.T) {
	svc := new(mocks.MockURLService)
	log := new(mocks.MockLogger)
	h := handler.NewHandler(svc, log)

	log.On("Info", "Handling shorten request", mock.Anything).Once()
	log.On("Warn", "Invalid shorten request", mock.Anything).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodPost, "/shorten", bytes.NewBufferString(`{"url":`)) // invalid
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	h.ShortenURL(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")

	svc.AssertNotCalled(t, "Shorten", mock.Anything, mock.Anything)
	log.AssertExpectations(t)
}

func TestShortenURL_ServiceError(t *testing.T) {
	svc := new(mocks.MockURLService)
	log := new(mocks.MockLogger)
	h := handler.NewHandler(svc, log)

	url := "https://fail.com"
	svc.On("Shorten", mock.Anything, url).Return("", errors.New("db down")).Once()

	log.On("Info", "Handling shorten request", mock.Anything).Once()
	log.On("Error", "Failed to shorten URL", mock.Anything, mock.Anything).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"url":"` + url + `"}`
	req, _ := http.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	h.ShortenURL(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to shorten URL")

	svc.AssertExpectations(t)
	log.AssertExpectations(t)
}

func TestResolveURL_Success(t *testing.T) {
	svc := new(mocks.MockURLService)
	log := new(mocks.MockLogger)
	h := handler.NewHandler(svc, log)

	slug := "abc123"
	originalURL := "https://go.dev"

	// Настраиваем мок
	svc.On("Resolve", mock.Anything, slug).Return(originalURL, nil).Once()
	log.On("Info", "Handling resolve request", mock.Anything).Maybe()
	log.On("Info", "Redirecting to original URL", mock.Anything).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Имитируем GET /:slug
	req, _ := http.NewRequest(http.MethodGet, "/abc123", nil)
	c.Params = gin.Params{gin.Param{Key: "slug", Value: slug}}
	c.Request = req

	h.ResolveURL(c)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, originalURL, w.Header().Get("Location"))

	svc.AssertExpectations(t)
	log.AssertExpectations(t)
}

func TestResolveURL_ServiceError(t *testing.T) {
	svc := new(mocks.MockURLService)
	log := new(mocks.MockLogger)
	h := handler.NewHandler(svc, log)

	slug := "fail"
	testErr := errors.New("db error")

	svc.On("Resolve", mock.Anything, slug).Return("", testErr).Once()

	log.On("Info", "Handling resolve request", mock.Anything).Once()
	log.On("Error", "Failed to resolve URL", testErr, mock.Anything).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/fail", nil)
	c.Params = gin.Params{gin.Param{Key: "slug", Value: slug}}
	c.Request = req

	h.ResolveURL(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to resolve URL")

	svc.AssertExpectations(t)
	log.AssertExpectations(t)
}

func TestResolveURL_NotFound(t *testing.T) {
	svc := new(mocks.MockURLService)
	log := new(mocks.MockLogger)
	h := handler.NewHandler(svc, log)

	slug := "notfound"
	// Метод Resolve вернул пустую строку (типа slug не найден)
	svc.On("Resolve", mock.Anything, slug).Return("", nil).Once()

	log.On("Info", "Handling resolve request", mock.Anything).Once()
	log.On("Warn", "Slug not found", mock.Anything).Once()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/notfound", nil)
	c.Params = gin.Params{gin.Param{Key: "slug", Value: slug}}
	c.Request = req

	h.ResolveURL(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Slug not found")

	svc.AssertExpectations(t)
	log.AssertExpectations(t)
}
