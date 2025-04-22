package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Handler {
	return &Handler{
		rdb: rdb,
	}
}

func (h *Handler) CreateTicket(c echo.Context) error {
	// ticket 생성 로직
	return c.String(200, "ticket created")
}

func RegisterRoute(e *echo.Echo, handler *Handler) {
	// ticket 등록
	e.POST("/ticket", handler.CreateTicket)
}
