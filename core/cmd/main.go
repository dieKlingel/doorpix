package main

import (
	"context"
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
	withSystem(app, &system)
	withHTTP(app, &system)
	withSIPPhone(app, &system)

	ctx := context.Background()
	app.Exec(ctx)
}

func withSystem(app *core.App, system *doorpix.System) {
	app.RegisterService(&core.SystemService{
		System: *system,
	})
}

func withHTTP(app *core.App, system *doorpix.System) {
	if system.Config.HTTP.Enabled {
		slog.Info("http is enabled")
		app.RegisterService(&core.HTTPService{
			System: *system,
		})
	}
}

func withSIPPhone(app *core.App, system *doorpix.System) {
	if system.Config.SIPPhone.Enabled {
		slog.Info("sip-phone is enabled")
		app.RegisterService(&core.PJSIPService{
			System: *system,
		})
	}
}
