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

	bus := core.NewBus()
	app := core.NewAppWithConfig(config, bus)
	withSystem(app, bus, config)
	withHTTP(app, bus, config)
	withSIPPhone(app, bus, config)
	withMQTT(app, bus, config)
	withRPC(app, bus, config)

	ctx := context.Background()
	app.Exec(ctx)
}

func withSystem(app *core.App, bus *core.Bus, config doorpix.Config) {
	app.RegisterService(&core.SystemService{
		Bus:    bus,
		Config: config,
	})
}

func withHTTP(app *core.App, _ *core.Bus, config doorpix.Config) {
	if config.HTTP.Enabled {
		slog.Info("http is enabled")
		app.RegisterService(&core.HTTPService{
			Config: config,
		})
	}
}

func withSIPPhone(app *core.App, bus *core.Bus, config doorpix.Config) {
	if config.SIPPhone.Enabled {
		slog.Info("sip-phone is enabled")
		app.RegisterService(&core.PJSIPService{
			Config: config,
			Emit:   bus.Write,
		})
	}
}

func withMQTT(app *core.App, bus *core.Bus, config doorpix.Config) {
	if config.MQTT.Enabled {
		slog.Info("mqtt is enabled")
		app.RegisterService(&core.MQTTService{
			Config: config,
			Emit:   bus.Write,
		})
	}
}

func withRPC(app *core.App, bus *core.Bus, config doorpix.Config) {
	if config.RPC.Enabled {
		slog.Info("rpc is enabled")
		app.RegisterService(&core.RPCService{
			Config:  config,
			Emitter: bus.Write,
		})
	}
}
