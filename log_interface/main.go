package main

import (
	"context"
	"io"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

/*
	지원해야 할 logging library
		1. zerolog
		2. zap
		3. log/slog
		4. logrus
		5. log
*/

type Logger interface {
	WithContext(ctx context.Context) context.Context
}

type logger struct {
	l Logger
}

// zap adapter
type ZapAdapter struct {
	*zap.Logger
}

type ctxKey struct{}

func (zp ZapAdapter) WithContext(ctx context.Context) context.Context {
	if _, ok := ctx.Value(ctxKey{}).(*Logger); !ok {
		//
		return ctx
	}

	return context.WithValue(ctx, ctxKey{}, &zp)
}

// handler
func handle(ctx context.Context) error {
	_ = ctx
	return nil
}

func main() {
	var w io.Writer
	var ctx = context.Background()
	zeroLogger := zerolog.New(w)
	_ = zeroLogger

	zapLogger, _ := zap.NewProduction()
	zp := ZapAdapter{zapLogger}

	logger := logger{l: zp}
	ctx = logger.l.WithContext(ctx)
	handle(ctx)
}
