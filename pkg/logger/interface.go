package logger

// Logger — интерфейс логгирования, позволяющий легко заменять реализацию в тестах.
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Fatal(msg string, err error, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
}
