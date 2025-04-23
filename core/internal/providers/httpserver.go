package providers

import (
	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"go.uber.org/fx"
)

func NewHTTPServer(lifecycle fx.Lifecycle, config doorpix.Config) *core.HTTPServer {
	if !config.HTTP.Enabled {
		return nil
	}

	server := core.NewHTTPServer(
		core.HTTPServerProps{
			Port:                    config.HTTP.Port,
			VideoStreamCameraDevice: config.Camera.Device,
		},
	)

	lifecycle.Append(
		fx.StartStopHook(server.Start, server.Stop),
	)

	return server
}
