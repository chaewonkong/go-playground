package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

var (
	_logger *Logger = newNop()
	mu      sync.RWMutex
)

func newNop() *Logger {
	return &Logger{s: slog.New(slog.NewTextHandler(io.Discard, nil))}
}

// G global logger 반환
func G() *Logger {
	mu.RLock()
	defer mu.RUnlock()

	return _logger
}

// Logger logger 구조체
type Logger struct {
	s *slog.Logger
}

// Debug debug level log with context
func (l *Logger) Debug(ctx context.Context, msg string, fields ...any) {
	l.s.DebugContext(ctx, msg, fields...)
}

// Info info level log with context
func (l *Logger) Info(ctx context.Context, msg string, fields ...any) {
	l.s.InfoContext(ctx, msg, fields...)
}

// Warn warn level log with context
func (l *Logger) Warn(ctx context.Context, msg string, fields ...any) {
	l.s.WarnContext(ctx, msg, fields...)
}

// Error error level log with context
func (l *Logger) Error(ctx context.Context, msg string, fields ...any) {
	l.s.ErrorContext(ctx, msg, fields...)
}

// InitGlobalLogger global logger 초기화
func InitGlobalLogger(optFns ...HandlerOptionFunc) {
	mu.Lock()
	defer mu.Unlock()

	handler := NewHandler(slog.NewJSONHandler(os.Stdout, nil), optFns...)
	l := slog.New(handler)

	_logger = &Logger{s: l}
}
