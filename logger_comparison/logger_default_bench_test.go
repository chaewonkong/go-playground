package loggercomparison

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

func BenchmarkZap(b *testing.B) {
	logger := newZapLogger(false)
	defer logger.Sync()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("Hello, world!", zap.String("key", "value"))
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

func BenchmarkLogrus(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

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
