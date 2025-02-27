package loggercomparison

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func BenchmarkZap(b *testing.B) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"discard"}
	logger, _ := cfg.Build()
	defer logger.Sync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Hello, world!", zap.String("key", "value"), zap.Time("time", time.Now()))
	}
}

func BenchmarkZeroLog(b *testing.B) {
	logger := zerolog.New(io.Discard).With().Timestamp().Logger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Str("key", "value").Msg("Hello, world!")
	}
}

func BenchmarkSlog(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Hello, world!", "key", "value")
	}
}

func BenchmarkLogrus(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(io.Discard)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(
			logrus.Fields{
				"key":  "value",
				"time": time.Now(),
			},
		).Info("Hello, world!")
	}
}
