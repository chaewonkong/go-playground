package benchmarks_test

import (
	"context"
	"os"
	"testing"

	"github.com/yuseferi/zax/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var someFields = []zap.Field{
	zap.String("field1", "value1"),
	zap.String("field2", "value2"),
	zap.Int("field3", 2),
}

func newZapLogger(isBuffer bool) (logger *zap.Logger, cleanup func()) {
	cleanupFns := []func(){}

	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	ec.TimeKey = "time"
	ec.CallerKey = ""
	ec.MessageKey = "message"
	enc := zapcore.NewJSONEncoder(ec)

	var ws zapcore.WriteSyncer = os.Stdout
	if isBuffer {
		bws := &zapcore.BufferedWriteSyncer{WS: ws, Size: 256 * 1024}
		cleanupFns = append(cleanupFns, func() { _ = bws.Stop() })
		ws = bws
	}

	logger = zap.New(zapcore.NewCore(
		enc,
		ws,
		zap.InfoLevel,
	))
	cleanupFns = append(cleanupFns, func() { _ = logger.Sync() })

	return logger, func() {
		for _, fn := range cleanupFns {
			fn()
		}
	}
}

func BenchmarkZapWithZax(b *testing.B) {
	logger, cleanup := newZapLogger(true)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := context.Background()
			ctx = zax.Set(ctx, someFields)
			logger.With(zax.Get(ctx)...).Info("hello, world")
		}
	})
}

func BenchmarkZap(b *testing.B) {
	logger, cleanup := newZapLogger(true)
	defer cleanup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.With(someFields...).Info("hello, world")
		}
	})
}
