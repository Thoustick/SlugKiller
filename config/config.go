package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config хранит все ключевые настройки приложения.
// Внутри полей используем типы, удобные для использования (int, time.Duration и т.п.)
type Config struct {
	HTTPAddr    string // Адрес, на котором слушает Gin
	StorageType string // "postgres" или "memory"
	DatabaseURL string // Postgres DSN
	RedisHost   string
	RedisPass   string
	RedisDB     int
	DBTimeout   time.Duration // Для context.WithTimeout
	SlugLength  int           // Длина короткой ссылки
	MaxAttempts int           // Число попыток при генерации
	CacheTTL    time.Duration // Для кеша в Redis
	LogLevel    string        // Уровень логирования
}

// Load создает экземпляр Config, считав значения из окружения.
// Если нужная переменная окружения пустая, используем fallback по умолчанию.
func Load() *Config {
	// Попытаемся загрузить .env, игнорируем ошибку, если файла нет
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.HTTPAddr = getEnv("HTTP_ADDR", ":8080")
	cfg.StorageType = getEnv("STORAGE_TYPE", "postgres")
	cfg.DatabaseURL = getEnv("DATABASE_URL", "postgres://user:pass@localhost:5432/shortener?sslmode=disable")

	cfg.RedisHost = getEnv("REDIS_HOST", "localhost:6379")
	cfg.RedisPass = getEnv("REDIS_PASSWORD", "")
	cfg.RedisDB = getEnvAsInt("REDIS_DB", 0)

	cfg.DBTimeout = getEnvAsDurationSeconds("DB_TIMEOUT_SECONDS", 15)

	cfg.SlugLength = getEnvAsInt("SLUG_LENGTH", 10)
	cfg.MaxAttempts = getEnvAsInt("MAX_ATTEMPTS", 5)

	// TTL в часах
	hours := getEnvAsInt("CACHE_TTL_HOURS", 0)
	cfg.CacheTTL = time.Duration(hours) * time.Hour

	cfg.LogLevel = getEnv("LOG_LEVEL", "info")
	return cfg
}

// getEnv возвращает значение переменной окружения или fallback, если она не определена
func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

// getEnvAsInt аналогично, но возвращает int
func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if valInt, err := strconv.Atoi(valStr); err == nil {
		return valInt
	}
	return fallback
}

// getEnvAsDurationSeconds возвращает Duration, считанную из переменной окружения
// как число секунд, иначе fallback (в сек)
func getEnvAsDurationSeconds(key string, fallback int) time.Duration {
	valStr := getEnv(key, "")
	if valInt, err := strconv.Atoi(valStr); err == nil {
		return time.Duration(valInt) * time.Second
	}
	return time.Duration(fallback) * time.Second
}
