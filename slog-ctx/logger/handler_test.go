package logger_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"

	"slog-ctx/logger"

	"github.com/stretchr/testify/assert"
)

func TestSlogCustomHandlerWithContext_defaultContextFields(t *testing.T) {
	type contextField struct {
		label string
		key   any
		value string
	}
	type versionKey struct{}
	type requestIDKey struct{}

	for _, tc := range []struct {
		name         string
		ctxField     []contextField
		ctxWithValue bool
	}{
		{
			name:         "handler.defaultContextFields 생성 시점에 추가한 ctx value를 로깅",
			ctxField:     []contextField{{label: "version", key: versionKey{}, value: "v1.0.11"}, {label: "request_id", key: requestIDKey{}, value: "abcd-1234"}},
			ctxWithValue: true,
		},
		{
			name:         "logger.DefaultContextFields ctx key가 context에 없으면 로깅하지 않음",
			ctxField:     []contextField{{label: "version", key: versionKey{}, value: "v1.0.11"}, {label: "request_id", key: requestIDKey{}, value: "abcd-1234"}},
			ctxWithValue: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			opts := []logger.HandlerOptionFunc{}
			for _, field := range tc.ctxField {
				opts = append(opts, logger.WithContextField(field.label, field.key))
			}

			h := logger.NewHandler(
				slog.NewJSONHandler(buf, nil),
				opts...,
			)
			l := slog.New(h)

			ctx := context.Background()

			// tc.ctxWithValue가 true인 경우에만 ctx에 key, value추가
			if tc.ctxWithValue {
				for _, field := range tc.ctxField {
					ctx = context.WithValue(ctx, field.key, field.value)
				}
			}
			l.InfoContext(ctx, "Hello", slog.String("name", "leon"))

			s := buf.String()
			for _, field := range tc.ctxField {
				fieldStr := fmt.Sprintf("%q:%q", field.label, field.value)
				if tc.ctxWithValue {
					assert.Contains(t, s, fieldStr)
				} else {
					assert.NotContains(t, s, fieldStr)
				}
			}
		})
	}
}

func TestSlogCustomHandlerWithContext_Prepend(t *testing.T) {
	t.Run("Prepend된 필드가 로깅하는 필드보다 먼저 로깅됨", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		l := slog.New(logger.NewHandler(slog.NewJSONHandler(buf, nil)))

		ctx := context.Background()
		ctx = logger.Prepend(ctx, slog.String("version", "v1.0.1"), slog.String("request_id", "1234"))

		l.InfoContext(ctx, "Hello", slog.String("name", "leon"))

		assert.Contains(t, buf.String(), `"version":"v1.0.1","request_id":"1234","name":"leon"`)

		buf.Reset() // clear buffer

		// 한번 더 prepend
		ctx = logger.Prepend(ctx, slog.String("version", "v1.0.2"))
		l.InfoContext(ctx, "Hello", slog.String("name", "leon"))

		assert.Contains(t, buf.String(), `"version":"v1.0.1","request_id":"1234","version":"v1.0.2","name":"leon"`)
	})
	t.Run("concurrency에서 race condition 체크 with go test -race", func(t *testing.T) {
		l := slog.New(logger.NewHandler(slog.NewJSONHandler(io.Discard, nil)))
		base := make([]slog.Attr, 0, 1000) // 큰 capacity 확보
		ctx := context.WithValue(context.Background(), logger.PrependKey{}, base)

		var wg sync.WaitGroup
		for i := range 1000 {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				attr := slog.String("pool_id", fmt.Sprintf("pool-%d", i))
				ctx2 := logger.Prepend(ctx, attr)
				l.InfoContext(ctx2, "childlog")
			}(i)
		}
		wg.Wait()
	})
}

func TestSlogCustomHandlerWithContext_WithGOA(t *testing.T) {
	// GOA: Group Or Attributes (그룹 또는 속성)에 대한 테스트
	t.Run("WithGroup, WithAttrs가 정상 동작함", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		l := slog.New(logger.NewHandler(slog.NewJSONHandler(buf, nil)))

		ctx := context.Background()
		ctx = logger.Prepend(ctx, slog.String("version", "v1.0.1"), slog.String("request_id", "1234"))

		logger := l.WithGroup("response").With(slog.String("status", "ok")).With(slog.Int("took", 123)).WithGroup("detail").With(slog.String("name", "leon"))
		logger.InfoContext(ctx, "Hello")

		s := buf.String()

		assert.Contains(t, s, `"version":"v1.0.1","request_id":"1234","response":{"status":"ok","took":123,"detail":{"name":"leon"}}}`)
	})
}
