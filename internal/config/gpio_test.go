package config_test

import (
	"testing"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestGPIO_UnmarshalYAML(t *testing.T) {
	t.Run("should parse int", func(t *testing.T) {
		input := "" +
			"enabled: true\r\n" +
			"chip: gpio0\r\n" +
			"inputs: [1, 2, 3]\r\n"

		gpio := &config.GPIO{}
		err := yaml.Unmarshal([]byte(input), gpio)

		assert.NoError(t, err)
		assert.Equal(t, true, gpio.Enabled)
		assert.Equal(t, "gpio0", gpio.Chip)
		assert.ElementsMatch(t, []int{1, 2, 3}, gpio.Inputs)
	})
}
