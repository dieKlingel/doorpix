package providers

import (
	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"go.uber.org/fx"
)

func NewGRPCServer(lifecycle fx.Lifecycle, eventemitter *eventemitter.EventEmitter, config doorpix.Config) *core.GRPCServer {
	if !config.RPC.Enabled {
		return nil
	}

	server := core.NewGRPCServer(eventemitter, core.GRPCServerProps{
		Address: config.RPC.Host,
		Port:    config.RPC.Port,
	})
	lifecycle.Append(fx.StartStopHook(server.Start, server.Stop))

	return server
}
