package logger_test

import (
	"log/slog"
	"testing"
	"time"

	. "slog-ctx/logger"

	"github.com/stretchr/testify/assert"
)

func TestConverTContextValueToSlogAttr(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		value    any
		expected slog.Attr
	}{
		{"string", "label", "hello", slog.String("label", "hello")},
		{"int", "count", 42, slog.Int("count", 42)},
		{"int64", "id", int64(123456), slog.Int64("id", 123456)},
		{"bool", "enabled", true, slog.Bool("enabled", true)},
		{"float64", "pi", 3.14, slog.Float64("pi", 3.14)},
		{"duration", "timeout", time.Second, slog.Duration("timeout", time.Second)},
		{"time", "now", time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC), slog.Time("now", time.Date(2023, 3, 1, 12, 0, 0, 0, time.UTC))},
		{"uint64", "big", uint64(999), slog.Uint64("big", 999)},
		{"fallback", "unknown", []int{1, 2}, slog.Any("unknown", []int{1, 2})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConverTContextValueToSlogAttr(tt.label, tt.value)
			assert.Equal(t, tt.expected, got)
		})
	}
}
