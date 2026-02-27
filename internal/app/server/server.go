package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/dieklingel/doorpix/internal/media/camera"
	"github.com/dieklingel/doorpix/internal/transport/http"
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

type Server struct {
	cameraDriver camera.Driver
	userAgent    *sip.UserAgent
	httpServer   *http.Server
}

func New(cfg *config.Config) *Server {
	cameraDriver := must(camera.NewGstDriver(`
		autovideosrc ! video/x-raw,width=800,height=600,framerate=20/1 ! tee name=tee
			tee. ! queue ! valve name=valve-http-camera ! jpegenc ! appsink name=appsink-http-camera
			tee. ! queue ! valve name=valve-sip-camera ! videoscale ! videoconvert ! video/x-raw,format=I420,width=720,height=480 ! appsink name=appsink-sip-camera
	`))

	var userAgent *sip.UserAgent = nil
	if cfg.SIP.Enabled {
		userAgent = sip.NewUserAgent(sip.UserAgentProps{
			Username: cfg.SIP.Username,
			Password: cfg.SIP.Password,
			Realm:    cfg.SIP.Realm,
			Domain:   cfg.SIP.Server,
			Webcam:   must(camera.NewWebcam("sip-camera", cameraDriver)),
		})
	}

	var httpServer *http.Server = nil
	if cfg.HTTP.Enabled {
		httpServer = http.NewServer(http.ServerProps{
			Webcam:    must(camera.NewWebcam("http-camera", cameraDriver)),
			UserAgent: userAgent,
			Port:      &cfg.HTTP.Port,
		})
	}

	return &Server{
		cameraDriver: cameraDriver,
		userAgent:    userAgent,
		httpServer:   httpServer,
	}
}

func (s *Server) Exec() {
	// Startup
	slog.Info("app server: starting up doorpix")
	if s.httpServer != nil {
		go func() {
			err := s.httpServer.Serve()
			if err != nil {
				slog.Error(err.Error())
			}
		}()
	}

	if s.userAgent != nil {
		go func() {
			err := s.userAgent.Serve()
			if err != nil {
				slog.Error(err.Error())
			}
		}()
	}

	// Catch Signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	slog.Info("app server: shutting down doorpix")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Shutdown
	if s.httpServer != nil {
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}

	if s.userAgent != nil {
		err := s.userAgent.Shutdown(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
