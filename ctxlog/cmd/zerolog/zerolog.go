package main

import (
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Service struct{}

func (s *Service) Hello(c echo.Context, name string) string {
	ctx := c.Request().Context()

	// zerolog.Ctx(ctx)는 ctx에 저장된 *zerolog.Logger를 반환함.
	// {"level":"info","request_id":"706e54a9-6aca-4032-89ff-e39af38eeaba","version":"v1.0.1","name":"leon","time":"2025-03-24T18:44:39+09:00","message":"Hello"}
	zerolog.Ctx(ctx).Info().Str("name", name).Msg("Hello")

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

		// zerolog로 context에 필드 추가
		logger := log.
			Output(os.Stderr). // stderr로 로그 출력
			With().
			Str("request_id", reqID).
			Str("version", "v1.0.1").
			Logger()

		// returns context.WithValue(ctx, ctxKey{}, &logger)
		ctx := logger.WithContext(req.Context())

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

	svc := &Service{}
	h := &Handler{svc: svc}
	e := NewServer(h)
	e.Logger.Fatal(e.Start(":8080"))
}
