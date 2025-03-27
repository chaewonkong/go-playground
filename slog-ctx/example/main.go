package main

import (
	"context"
	"log/slog"
	"net/http"

	"slog-ctx/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var version = "v1.0.23"

type requestIDKey struct{}
type versionKey struct{}

const (
	requestIDKeyStr = "request_id"
	versionKeyStr   = "version"
)

type Service struct{}

func (s *Service) Hello(ctx context.Context, name string) string {
	logger.G().Info(ctx, "Hello", slog.String("name", name))
	// {"time":"2025-03-27T18:20:33.445328352+09:00","level":"INFO","msg":"Hello","name":"leon","request_id":"b09f47c6-9fba-436a-8511-903a3ee2d610","version":"v1.0.23","pool_id":1638,"name":"leon"}

	return "Hello, " + name
}

type Handler struct {
	svc *Service
}

func (h *Handler) Hello(c echo.Context) error {
	name := c.Param("name")
	ctx := c.Request().Context()

	// pool_id를 로깅 context에 추가
	ctx = logger.Prepend(ctx, slog.Int("pool_id", 1638))

	return c.String(http.StatusOK, h.svc.Hello(ctx, name))
}

// LoggerMiddleware requestID, version을 로깅 context에 추가하는 미들웨어
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		reqID := uuid.New().String()

		// requestID, version을 context에 추가
		// Prepend: 로깅 시 실제 로깅하려는 메시지 필드보다 앞에 추가한 필드가 로깅되도록 설정. Append는 반대.
		ctx := context.WithValue(req.Context(), requestIDKey{}, reqID)
		ctx = context.WithValue(ctx, versionKey{}, version)

		c.SetRequest(req.WithContext(ctx))

		return next(c)
	}
}

func NewServer(h *Handler) *echo.Echo {
	e := echo.New()
	e.GET("/hello/:name", h.Hello)

	e.Use(LoggerMiddleware, middleware.RequestID())

	return e
}

func main() {
	// logger setting
	logger.InitGlobalLogger(
		logger.WithContextField(requestIDKeyStr, requestIDKey{}),
		logger.WithContextField(versionKeyStr, versionKey{}),
	)

	svc := &Service{}
	h := &Handler{svc: svc}
	e := NewServer(h)
	e.Logger.Fatal(e.Start(":8080"))
}
