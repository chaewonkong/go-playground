package server

import (
	"context"
	"gracefulshutdown/service"
	"gracefulshutdown/service/apiserver/handler"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Port  string `env:"SERVER_PORT,default:8080"`
	Redis struct {
		Host string `env:"REDIS_HOST,default:localhost"`
		Port string `env:"REDIS_PORT,default:6379"`
	}
}

func NewApp(cfg *Config) (service.Application, error) {
	// redis 연결
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Host + ":" + cfg.Redis.Port,
	})

	// redis 연결 확인
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	// HTTP server 생성
	e := echo.New()

	// handler 생성
	h := handler.New(rdb)

	// route 등록
	handler.RegisterRoute(e, h)

	return &App{
		config: cfg,
		echo:   e,
	}, nil
}

type App struct {
	config *Config
	echo   *echo.Echo
}

func (s *App) Run() error {
	return s.echo.Start(":" + s.config.Port)
}

func (s *App) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
