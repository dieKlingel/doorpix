package providers

import (
	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"go.uber.org/fx"
)

func NewGPIOController(lifecycle fx.Lifecycle, eventemitter *eventemitter.EventEmitter, config doorpix.Config) *core.GPIOController {
	if !config.GPIO.Enabled {
		return nil
	}

	gpioController := core.NewGPIOController(
		eventemitter,
		core.GPIOControllerProps{
			Pins: config.GPIO.Pins,
			Chip: config.GPIO.Chip,
		},
	)

	lifecycle.Append(fx.StartStopHook(gpioController.Start, gpioController.Stop))
	return gpioController
}
