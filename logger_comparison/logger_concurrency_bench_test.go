package loggercomparison

import (
	"io"
	"log/slog"
	"math/rand"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func BenchmarkZapConcurrent(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"discard"}
	logger, _ := cfg.Build()
	defer logger.Sync()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Concurrent log", zap.String("key", "value"), zap.Int("id", rand.Int()), zap.Time("time", time.Now()))
		}
	})
}

func BenchmarkZeroLogConcurrent(b *testing.B) {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Str("key", "value").Int("id", rand.Int()).Msg("Concurrent log")
		}
	})
}

func BenchmarkSlogConcurrent(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Concurrent log", slog.String("key", "value"), slog.Int("id", rand.Int()))
		}
	})
}

func BenchmarkLogrusConcurrent(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(io.Discard)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(logrus.Fields{
				"key":  "value",
				"id":   rand.Int(),
				"time": time.Now(),
			}).Info("Concurrent log")
		}
	})
}
