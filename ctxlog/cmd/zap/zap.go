package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yuseferi/zax/v2"
	"go.uber.org/zap"
)

type Service struct{}

func (s *Service) Hello(c echo.Context, name string) string {
	ctx := c.Request().Context()

	// 로깅할 때마다 zax.Get으로 ctx에 저장된 필드를 가져와야 함
	// {"level":"info","ts":1742808992.2126977,"caller":"zap/zap.go:19","msg":"Hello","request_id":"3cc27a4e-5c86-4750-8e6c-4f29d0381bc5","version":"v1.0.1","name":"leon"}
	zap.L().With(zax.Get(ctx)...).Info("Hello", zap.String("name", name))

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

		// zax로 context에 필드 추가
		ctx := zax.Set(req.Context(), []zap.Field{zap.String("request_id", reqID), zap.String("version", "v1.0.1")})

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
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	logger := zap.Must(cfg.Build())
	defer func() { _ = logger.Sync() }()

	reset := zap.ReplaceGlobals(logger)
	defer reset()

	svc := &Service{}
	h := &Handler{svc: svc}
	e := NewServer(h)
	e.Logger.Fatal(e.Start(":8080"))
}
