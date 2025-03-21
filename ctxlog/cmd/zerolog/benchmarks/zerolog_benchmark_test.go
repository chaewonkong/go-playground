package benchmarks

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func BenchmarkZerolog(b *testing.B) {
	ctx := context.Background()
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	ctx = logger.WithContext(ctx)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger := zerolog.Ctx(ctx)
			logger.UpdateContext(
				func(c zerolog.Context) zerolog.Context {
					uuid := uuid.New().String()
					return c.Str("uuid", uuid)
				},
			)

			logger.Info().Msg("hello, world")
		}

	})
}
