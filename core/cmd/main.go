package main

import (
	"log/slog"
	"os"

	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/env"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/providers"
	"go.uber.org/fx"
)

type App struct {
	fx.In

	App        *core.App        `optional:"false"`
	HTTPServer *core.HTTPServer `optional:"true"`
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
			providers.NewApp,
		),
		fx.Invoke(func(app App) {
			// do nothing, the core.App is started by the fx.Lifecycle
		}),
	)

	app.Run()
}
