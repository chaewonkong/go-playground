package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	e := echo.New()

	loggerMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// Request ID from Echo's built-in header
			requestID := res.Header().Get(echo.HeaderXRequestID)

			// 새 context와 logger 생성
			logger := log.With().
				Str("request_id", requestID).
				Logger()

			ctx := logger.WithContext(req.Context())

			// context 교체
			req = req.WithContext(ctx)
			c.SetRequest(req)

			return next(c)
		}
	}

	e.Use(middleware.RequestID(), loggerMiddleware)

	e.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()
		logger := zerolog.Ctx(ctx)

		logger.Info().Msg("hello, world") //{"level":"info","request_id":"AHNqjLRPcHkRQQJjaGVxFgTYcXxEuKKw","time":"2025-03-21T15:01:05+09:00","message":"hello, world"}
		logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Bool("production", false).Str("foo", "bar")
		})

		logger.Info().Msg("bye, world")
		return c.String(http.StatusOK, "hello, world")
	})

	e.GET("/updated", func(c echo.Context) error {
		ctx := c.Request().Context()

		logger := zerolog.Ctx(ctx)

		// log fields from ctx
		logger.Info().Msg("hello, world")

		return c.String(http.StatusOK, "updated")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
