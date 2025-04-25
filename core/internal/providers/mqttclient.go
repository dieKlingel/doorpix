package providers

import (
	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"go.uber.org/fx"
)

func NewMQTTClient(lifecycle fx.Lifecycle, config doorpix.Config, eventemitter *eventemitter.EventEmitter) *core.MQTTClient {
	if !config.MQTT.Enabled {
		return nil
	}

	client := core.NewMQTTClient(
		eventemitter,
		core.MQTTClientProps{
			ClientId:      config.MQTT.ClientId,
			Host:          config.MQTT.Host,
			Port:          config.MQTT.Port,
			Protocol:      config.MQTT.Protocol,
			Username:      config.MQTT.Username,
			Password:      config.MQTT.Password,
			Subscriptions: config.MQTT.Subscriptions,
		},
	)

	lifecycle.Append(fx.StartStopHook(client.Start, client.Stop))
	return client
}
