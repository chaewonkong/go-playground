package main

import (
	"gracefulshutdown/service"
	"gracefulshutdown/service/apiserver/server"
	"log/slog"
	"os"
)

type c struct{}

func main() {
	cfg := &server.Config{}

	err := service.RunApplication(server.NewApp, cfg)
	if err != nil {
		slog.Default().Error("failed to run", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
