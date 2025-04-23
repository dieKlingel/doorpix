package providers

import (
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/service/httpsvc"
)

func NewHTTPServer(config doorpix.Config) *httpsvc.HTTPService {
	if !config.HTTP.Enabled {
		return nil
	}

	return httpsvc.New(
		httpsvc.HTTPServiceProps{
			Port:                    config.HTTP.Port,
			VideoStreamCameraDevice: config.Camera.Device,
		},
	)
}
