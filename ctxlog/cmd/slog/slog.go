package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	slogctx "github.com/veqryn/slog-context"
)

type Service struct{}

func (s *Service) Hello(c echo.Context, name string) string {
	ctx := c.Request().Context()

	// InfoContext, ErrContext 처럼 ctx와 함께 로깅해야만 추가된 필드가 로깅됨
	// {"time":"2025-03-24T18:24:46.876945824+09:00","level":"INFO","msg":"Hello","request_id":"50106a6c-0fda-41ce-b33a-e04f6934aebc","version":"v1.0.1","name":"leon"}
	slog.InfoContext(ctx, "Hello", slog.String("name", name))

	return "Hello, " + name
}

type Handler struct {
	svc *Service
}

func (h *Handler) Hello(c echo.Context) error {
	name := c.Param("name")

	return c.String(http.StatusOK, h.svc.Hello(c, name))
}

// LoggerMiddleware requestID, version을 로깅 context에 추가하는 미들웨어
func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		reqID := uuid.New().String()

		// requestID, version을 context에 추가
		// Prepend: 로깅 시 실제 로깅하려는 메시지 필드보다 앞에 추가한 필드가 로깅되도록 설정. Append는 반대.
		ctx := slogctx.Prepend(req.Context(), slog.String("request_id", reqID), slog.String("version", "v1.0.1"))

		c.SetRequest(req.WithContext(ctx))

		return next(c)
	}
}

func NewServer(h *Handler) *echo.Echo {
	e := echo.New()
	e.GET("/hello/:name", h.Hello)

	e.Use(LoggerMiddleware)

	return e
}

func main() {
	// logger setting
	logHandler := slogctx.NewHandler(slog.NewJSONHandler(os.Stderr, nil), nil)
	slog.SetDefault(slog.New(logHandler))

	svc := &Service{}
	h := &Handler{svc: svc}
	e := NewServer(h)
	e.Logger.Fatal(e.Start(":8080"))
}
