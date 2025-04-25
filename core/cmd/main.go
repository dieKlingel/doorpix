package main

import (
	"log/slog"
	"os"

	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/dieklingel/doorpix/core/internal/env"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/providers"
	"go.uber.org/fx"
)

type App struct {
	fx.In

	EventEmitter *eventemitter.EventEmitter `optional:"false"`
	App          *core.App                  `optional:"false"`
	HTTPServer   *core.HTTPServer           `optional:"true"`
	MQTTClient   *core.MQTTClient           `optional:"true"`
}

func (app *App) Start() error {
	err := app.EventEmitter.Emit("events/startup", map[string]any{})
	if err != nil {
		slog.Error("failed to emit start event", "error", err)
	}
	return nil
}

func (app *App) Stop() error {
	err := app.EventEmitter.Emit("events/shutdown", map[string]any{})
	if err != nil {
		slog.Error("failed to emit shutdown event", "error", err)
	}
	return nil
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: env.LogLevel(),
	}))
	slog.SetDefault(logger)

	app := fx.New(
		fx.Provide(
			eventemitter.New,
			providers.NewApplicationConfiguration,
			providers.NewHTTPServer,
			fx.Annotate(
				providers.NewMQTTClient,
				fx.As(new(actions.Publisher)),
			),
			providers.NewApp,
		),
		fx.Invoke(func(lifecyle fx.Lifecycle, app App) {
			lifecyle.Append(fx.StartStopHook(app.Start, app.Stop))
		}),
	)

	app.Run()
}
