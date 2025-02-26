package loggercomparison

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

func BenchmarkZap(b *testing.B) {
	logger := zap.NewExample()
	defer logger.Sync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Hello, world!", zap.String("key", "value"), zap.Time("time", time.Now()))
	}
}

func BenchmarkZeroLog(b *testing.B) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("key", "value").Msg("Hello, world!")
	}
}

func BenchmarkSlog(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Hello, world!", "key", "value")
	}
}
