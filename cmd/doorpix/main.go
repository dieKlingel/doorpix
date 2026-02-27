package main

import (
	"log/slog"
	"os"

	"github.com/dieklingel/doorpix/internal/app/server"
	"github.com/dieklingel/doorpix/internal/config"
)

func main() {
	cfg, err := config.NewBuilder().
		AddConfigFile("doorpix.yaml").
		AddConfigFile("~/doorpix.yaml").
		Build()

	if err != nil {
		slog.Error("failed to build config", "error", err)
		os.Exit(1)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	app := server.New(cfg)
	app.Exec()
}
