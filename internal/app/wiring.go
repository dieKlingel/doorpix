package app

import (
	"github.com/dieklingel/doorpix/internal/config"
	"github.com/dieklingel/doorpix/internal/media/camera"
	"github.com/dieklingel/doorpix/internal/transport/http"
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

func CreateCameraDriver(cfg *config.Config) camera.Driver {
	return must(camera.NewGstDriver(`
		autovideosrc ! video/x-raw,width=800,height=600,framerate=20/1 ! tee name=tee
			tee. ! queue ! valve name=valve-http-camera ! jpegenc ! appsink name=appsink-http-camera
	`))
}

func CreateUserAgent(cfg *config.Config) *sip.UserAgent {
	if !cfg.SIP.Enabled {
		return nil
	}

	return sip.NewUserAgent(sip.UserAgentProps{
		Username: cfg.SIP.Username,
		Password: cfg.SIP.Password,
		Realm:    cfg.SIP.Realm,
		Domain:   cfg.SIP.Server,
	})
}

func CreateHTTPServer(cfg *config.Config, driver camera.Driver, userAgent *sip.UserAgent) *http.Server {
	if !cfg.HTTP.Enabled {
		return nil
	}

	return http.NewServer(http.ServerProps{
		Webcam:    must(camera.NewWebcam("http-camera", driver)),
		UserAgent: userAgent,
		Port:      &cfg.HTTP.Port,
	})
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
