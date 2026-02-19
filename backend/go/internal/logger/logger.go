package logger

import "context"

// Logger is a structured, leveled logger interface used across the backend.
// It is intentionally minimal and compatible with zap's SugaredLogger.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})

	// With returns a child logger with additional context.
	With(keysAndValues ...interface{}) Logger

	// WithContext may attach request-scoped values such as request IDs.
	WithContext(ctx context.Context) Logger

	// Sync flushes any buffered log entries.
	Sync() error
}

