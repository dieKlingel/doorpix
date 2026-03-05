package config_test

import (
	"testing"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestHTTP_UnmarshalYAML(t *testing.T) {
	t.Run("with default", func(t *testing.T) {
		t.Run("should keep default", func(t *testing.T) {
			input := ""
			output := &config.HTTP{
				Enabled: true,
				Port:    1234,
			}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Enabled, true)
			assert.Equal(t, output.Port, 1234)
		})

		t.Run("should keep default when null", func(t *testing.T) {
			input := "" +
				"enabled: null\r\n" +
				"port: null"
			output := &config.HTTP{
				Enabled: true,
				Port:    1234,
			}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Enabled, true)
			assert.Equal(t, output.Port, 1234)
		})

		t.Run("should keep partial default", func(t *testing.T) {
			input := "port: 4321"
			output := &config.HTTP{
				Enabled: true,
				Port:    1234,
			}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Enabled, true)
			assert.Equal(t, output.Port, 4321)
		})

		t.Run("should overwrite all defaults", func(t *testing.T) {
			input := "" +
				"enabled: false\r\n" +
				"port: 2314"
			output := &config.HTTP{
				Enabled: true,
				Port:    1234,
			}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Enabled, false)
			assert.Equal(t, output.Port, 2314)
		})
	})

	t.Run("without defaults", func(t *testing.T) {
		t.Run("should fill with empty values", func(t *testing.T) {
			input := ""
			output := &config.HTTP{}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Enabled, false)
			assert.Equal(t, output.Port, 0)
		})

		t.Run("should set partial values", func(t *testing.T) {
			input := "port: 1234"
			output := &config.HTTP{}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Enabled, false)
			assert.Equal(t, output.Port, 1234)
		})

		t.Run("should set all values", func(t *testing.T) {
			input := "" +
				"enabled: true\r\n" +
				"port: 2468"
			output := &config.HTTP{}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Enabled, true)
			assert.Equal(t, output.Port, 2468)
		})
	})
}
