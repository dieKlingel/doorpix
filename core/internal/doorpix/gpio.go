package doorpix

import (
	"gopkg.in/yaml.v3"
)

type GPIOConfig struct {
	Enabled bool
	Chip    string
	Pins    map[int]string
}

func (gpioConfig *GPIOConfig) UnmarshalYAML(node *yaml.Node) error {
	rawGPIOConfig := struct {
		Enabled bool           `yaml:"enabled"`
		Chip    string         `yaml:"chip"`
		Pins    map[int]string `yaml:"pins"`
	}{
		Enabled: true,
		Pins:    make(map[int]string),
		Chip:    "gpiochip0",
	}

	if err := node.Decode(&rawGPIOConfig); err != nil {
		return err
	}

	gpioConfig.Enabled = rawGPIOConfig.Enabled
	gpioConfig.Chip = rawGPIOConfig.Chip
	gpioConfig.Pins = rawGPIOConfig.Pins
	return nil
}
