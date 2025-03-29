package logger

import (
	"os"
	"strings"
	"time"

	"github.com/Thoustick/SlugKiller/config"
	"github.com/rs/zerolog"
)

// ZeroLogger — конкретная реализация Logger с использованием zerolog.
type ZeroLogger struct {
	logger zerolog.Logger
}

// InitLogger создает логгер с конфигурацией
// func InitLogger(cfg *config.Config) Logger {
func InitLogger(cfg *config.Config) Logger {
	levelStr := strings.ToLower(cfg.LogLevel)
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}

	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	zlog := zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Logger()

	zlog.Info().Msg("Логгер успешно инициализирован")

	return &ZeroLogger{logger: zlog}
}

func (l *ZeroLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *ZeroLogger) Error(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Error().Fields(fields)
	if err != nil {
		event = event.Err(err)
	}
	event.Msg(msg)
}

func (l *ZeroLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

func (l *ZeroLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

func (l *ZeroLogger) Fatal(msg string, err error, fields map[string]interface{}) {
	event := l.logger.Fatal().Fields(fields)
	if err != nil {
		event = event.Err(err)
	}
	event.Msg(msg)
	os.Exit(1)
}
