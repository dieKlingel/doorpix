package doorpix_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/stretchr/testify/assert"
)

func TestCameraConfig_Error(t *testing.T) {
	t.Run("should return error on empty device", func(t *testing.T) {
		config := doorpix.CameraConfig{}

		err := config.Error()
		assert.Error(t, err)
	})

	t.Run("should return error on invalid device", func(t *testing.T) {
		config := doorpix.CameraConfig{
			Device: "invalid-device",
		}

		err := config.Error()
		assert.Error(t, err)
	})

	t.Run("should return nil on valid device", func(t *testing.T) {
		config := doorpix.CameraConfig{
			Device: "videotestsrc",
		}

		err := config.Error()
		assert.NoError(t, err)
	})
}
