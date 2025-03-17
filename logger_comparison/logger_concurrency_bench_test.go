package loggercomparison

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A Syncer is a spy for the Sync portion of zapcore.WriteSyncer.
type Syncer struct {
	err    error
	called bool
}

// SetError sets the error that the Sync method will return.
func (s *Syncer) SetError(err error) {
	s.err = err
}

// Sync records that it was called, then returns the user-supplied error (if
// any).
func (s *Syncer) Sync() error {
	s.called = true
	return s.err
}

// Called reports whether the Sync method was called.
func (s *Syncer) Called() bool {
	return s.called
}

// A Discarder sends all writes to io.Discard.
type Discarder struct{ Syncer }

// Write implements io.Writer.
func (d *Discarder) Write(b []byte) (int, error) {
	return io.Discard.Write(b)
}

func newZapLogger(isBuffer bool) *zap.Logger {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	ec.TimeKey = "time"
	ec.CallerKey = ""
	ec.MessageKey = "message"
	enc := zapcore.NewJSONEncoder(ec)

	var ws zapcore.WriteSyncer = os.Stdout
	if isBuffer {
		ws = &zapcore.BufferedWriteSyncer{WS: ws, Size: 256 * 1024}
	}

	return zap.New(zapcore.NewCore(
		enc,
		ws,
		zap.InfoLevel,
	))
}

func BenchmarkZapConcurrentWithBuffer(b *testing.B) {
	logger := newZapLogger(true)
	defer logger.Sync()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("concurrent-log", zap.String("key", "value"))
		}
	})
}

func BenchmarkZapConcurrent(b *testing.B) {

	logger := newZapLogger(false)
	defer logger.Sync()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("concurrent-log", zap.String("key", "value"))
		}
	})
}

// var mu sync.Mutex
// bwr := bufio.NewWriterSize(os.Stdout, 512*1024)

// logger := log.New(io.WriterFunc(func(p []byte) (n int, err error) {
//     mu.Lock()
//     defer mu.Unlock()
//     return bwr.Write(p)
// }), "", log.LstdFlags)

func BenchmarkZeroLogConcurrentWithBuffer(b *testing.B) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	bwr := bufio.NewWriterSize(os.Stdout, 256*1024)
	logger := zerolog.New(zerolog.SyncWriter(bwr)).With().Timestamp().Logger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Str("key", "value").Msg("concurrent-log")
		}
	})
}

func BenchmarkZeroLogConcurrent(b *testing.B) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info().Str("key", "value").Msg("concurrent-log")
		}
	})
}

type syncWriter struct {
	mu sync.Mutex
	wr io.Writer
}

func SyncWriter(w io.Writer) io.Writer {
	return &syncWriter{
		wr: w,
	}
}

func (s *syncWriter) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wr.Write(p)
}

func BenchmarkSlogConcurrentWithBuffer(b *testing.B) {
	bwr := bufio.NewWriterSize(os.Stdout, 256*1024)
	swr := SyncWriter(bwr)
	logger := slog.New(slog.NewJSONHandler(swr, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.Int64("time", a.Value.Time().UnixNano())
			case slog.MessageKey:
				return slog.String("message", a.Value.String())
			case slog.LevelKey:
				return slog.String("level", "info")
			}
			return a
		},
	}))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("concurrent-log", slog.String("key", "value"))
		}
	})
}

func BenchmarkSlogConcurrent(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.Int64("time", a.Value.Time().UnixNano())
			case slog.MessageKey:
				return slog.String("message", a.Value.String())
			case slog.LevelKey:
				return slog.String("level", "info")
			}
			return a
		},
	}))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("concurrent-log", slog.String("key", "value"))
		}
	})
}

func BenchmarkLogrusConcurrentWithBuffer(b *testing.B) {
	bwr := bufio.NewWriterSize(os.Stdout, 256*1024)
	swr := SyncWriter(bwr)
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:  "message",
			logrus.FieldKeyTime: "ts",
		}})
	logger.SetOutput(swr)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(logrus.Fields{
				"key":  "value",
				"time": time.Now().UnixNano(),
			}).Info("concurrent-log")
		}
	})
}

func BenchmarkLogrusConcurrent(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:  "message",
			logrus.FieldKeyTime: "ts",
		}})
	logger.SetOutput(os.Stdout)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(logrus.Fields{
				"key":  "value",
				"time": time.Now().UnixNano(),
			}).Info("concurrent-log")
		}
	})
}
