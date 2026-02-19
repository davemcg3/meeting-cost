package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
)

type zapLogger struct {
	base *zap.SugaredLogger
}

// NewZapLogger creates a new Zap-based Logger configured for the given
// environment ("development" or "production").
func NewZapLogger(env string) (Logger, error) {
	var cfg zap.Config
	if env == "production" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
	}

	// Always log to stdout so Docker can capture logs.
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	z, err := cfg.Build(zap.AddCaller())
	if err != nil {
		return nil, err
	}
	return &zapLogger{base: z.Sugar()}, nil
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.base.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.base.Infow(msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.base.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.base.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) With(keysAndValues ...interface{}) Logger {
	return &zapLogger{base: l.base.With(keysAndValues...)}
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}
	if reqID, ok := ctx.Value(ContextKeyRequestID).(string); ok && reqID != "" {
		return l.With("request_id", reqID)
	}
	return l
}

func (l *zapLogger) Sync() error {
	return l.base.Sync()
}

// ContextKey is a private type for context keys used by the logger middleware.
type ContextKey string

const (
	// ContextKeyRequestID is the context key under which the request ID is stored.
	ContextKeyRequestID ContextKey = "request_id"
)

// NewNopLogger returns a logger that discards all logs (useful in tests).
func NewNopLogger() Logger {
	return &zapLogger{base: zap.NewNop().Sugar()}
}

// DefaultLogger is a convenience for early bootstrap before DI container is built.
func DefaultLogger() Logger {
	l, err := NewZapLogger(os.Getenv("ENV"))
	if err != nil {
		return NewNopLogger()
	}
	return l
}

