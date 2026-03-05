package config_test

import (
	"testing"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestCamera_UnmarshalYAML(t *testing.T) {
	t.Run("with default", func(t *testing.T) {
		t.Run("should keep default", func(t *testing.T) {
			input := ""
			output := &config.Camera{
				Device: "default value",
			}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Device, "default value")
		})

		t.Run("should keep default when null", func(t *testing.T) {
			input := "device: null"
			output := &config.Camera{
				Device: "default value",
			}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Device, "default value")
		})

		t.Run("should overwrite all defaults", func(t *testing.T) {
			input := "device: hello"
			output := &config.Camera{
				Device: "default value",
			}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Device, "hello")
		})
	})

	t.Run("without defaults", func(t *testing.T) {
		t.Run("should fill with empty values", func(t *testing.T) {
			input := ""
			output := &config.Camera{}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Device, "")
		})

		t.Run("should set all values", func(t *testing.T) {
			input := "device: hello"
			output := &config.Camera{}

			err := yaml.Unmarshal([]byte(input), output)
			assert.NoError(t, err)
			assert.Equal(t, output.Device, "hello")
		})
	})
}
