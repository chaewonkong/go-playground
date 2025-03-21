package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/yuseferi/zax/v2"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()

	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	reset := zap.ReplaceGlobals(logger)
	defer reset()

	loggerMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()

			// Request ID from Echo's built-in header
			requestID := res.Header().Get(echo.HeaderXRequestID)

			ctx := req.Context()

			ctx = zax.Set(ctx, []zap.Field{zap.String("request_id", requestID)})

			// context 교체
			req = req.WithContext(ctx)
			c.SetRequest(req)

			return next(c)
		}
	}

	e.Use(middleware.RequestID(), loggerMiddleware)

	e.GET("/", func(c echo.Context) error {
		ctx := c.Request().Context()

		// log fields from ctx
		zap.L().With(zax.Get(ctx)...).Info("hello, world")

		// append more fields to ctx
		ctx = zax.Append(ctx, []zap.Field{zap.Bool("production", false), zap.String("foo", "bar")})
		zap.L().With(zax.Get(ctx)...).Info("bye, world")

		return c.String(http.StatusOK, "hello, world")
	})

	e.GET("/updated", func(c echo.Context) error {
		ctx := c.Request().Context()

		// log fields from ctx
		zap.L().With(zax.Get(ctx)...).Info("hello, world")

		return c.String(http.StatusOK, "updated")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
