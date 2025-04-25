package providers

import (
	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"go.uber.org/fx"
)

func NewSIPClient(lifecycle fx.Lifecycle, eventemitter *eventemitter.EventEmitter, config doorpix.Config) *core.SIPClient {
	if !config.SIPPhone.Enabled {
		return nil
	}

	sipClient := core.NewSIPClient(eventemitter, core.SIPClientProps{
		Username:    config.SIPPhone.Username,
		Password:    config.SIPPhone.Password,
		Realm:       config.SIPPhone.Realm,
		Server:      config.SIPPhone.Server,
		StunServers: config.SIPPhone.StunServers,
		VideoDevice: config.Camera.Device,
	})

	lifecycle.Append(fx.StartStopHook(sipClient.Start, sipClient.Stop))
	return sipClient
}
