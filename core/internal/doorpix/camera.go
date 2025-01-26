package doorpix

import (
	"fmt"
)

type CameraConfig struct {
	Device string `yaml:"device"`
}

func (c *CameraConfig) Error() error {
	if c.Device == "" {
		return fmt.Errorf("device is required")
	}

	return nil
}
