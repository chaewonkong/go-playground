package loggercomparison

import (
	"log/slog"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

func BenchmarkZapConcurrent(b *testing.B) {
	logger := zap.NewExample()
	defer logger.Sync()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Concurrent log", zap.String("key", "value"), zap.Int("id", rand.Int()), zap.Time("time", time.Now()))
		}
	})
}

func BenchmarkZeroLogConcurrent(b *testing.B) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Str("key", "value").Int("id", rand.Int()).Msg("Concurrent log")
		}
	})
}

func BenchmarkSlogConcurrent(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Concurrent log", slog.String("key", "value"), slog.Int("id", rand.Int()))
		}
	})
}
