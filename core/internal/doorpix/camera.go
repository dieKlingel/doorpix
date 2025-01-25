package doorpix

import (
	"fmt"

	"github.com/dieklingel/doorpix/core/internal/camera"
)

type CameraConfig struct {
	Device string `yaml:"device"`
}

func (c *CameraConfig) Error() error {
	if c.Device == "" {
		return fmt.Errorf("device is required")
	}

	_, err := camera.NewElement(c.Device)
	if err != nil {
		return err
	}

	return nil
}
