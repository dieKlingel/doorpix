package providers

import (
	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/service/httpsvc"
	"go.uber.org/fx"
)

type AppParams struct {
	fx.In

	Config       doorpix.Config
	EventEmitter *eventemitter.EventEmitter

	HTTPService *httpsvc.HTTPService `optional:"true"`
}

func NewApp(lifecycle fx.Lifecycle, params AppParams) *core.App {
	app := core.App{
		Config:       params.Config,
		EventEmitter: params.EventEmitter,
	}
	lifecycle.Append(
		fx.StartStopHook(app.Start, app.Stop),
	)

	return &app
}
