package main

import (
	"log/slog"
	"os"

	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/env"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: env.LogLevel(),
	}))
	slog.SetDefault(logger)

	config := doorpix.NewConfig()
	config.AddConfigPath(
		"/etc/doorpix/doorpix.yaml",
		"/etc/doorpix/config.yaml",
		"doorpix.yaml",
		"config.yaml",
	)
	if err := config.Read(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if err := config.Error(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	bus := core.NewEventEmitterWithConfig(config)
	system := doorpix.System{
		Config: config,
		Bus:    bus,
	}

	app := core.NewAppWithConfig(system)
	app.RegisterHandler(&core.SystemHandler{
		System: system,
	})
	app.RegisterHandler(&core.PJSIPPhone{
		System: system,
	})
	app.RegisterHandler(&core.HttpHandler{
		System: system,
	})

	app.Exec()
}
