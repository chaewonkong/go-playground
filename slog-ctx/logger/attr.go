package logger

import (
	"log/slog"
	"time"
)

// ConverTContextValueToSlogAttr context value를 slog.Attr로 변환
func ConverTContextValueToSlogAttr(label string, val any) slog.Attr {
	switch v := val.(type) {
	case string:
		return slog.String(label, v)
	case int:
		return slog.Int(label, v)
	case int64:
		return slog.Int64(label, v)
	case bool:
		return slog.Bool(label, v)
	case float64:
		return slog.Float64(label, v)
	case time.Duration:
		return slog.Duration(label, v)
	case time.Time:
		return slog.Time(label, v)
	case uint64:
		return slog.Uint64(label, v)
	default:
		return slog.Any(label, v)
	}
}
