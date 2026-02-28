package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/dieklingel/doorpix/internal/doorpi/system"
	"github.com/dieklingel/doorpix/internal/media/camera"
	"github.com/dieklingel/doorpix/internal/oplog"
	"github.com/dieklingel/doorpix/internal/transport/http"
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

type Server struct {
	cameraDriver camera.Driver
	userAgent    *sip.UserAgent
	httpServer   *http.Server
	shellService *system.ShellService
}

func New(cfg *config.Config) *Server {
	cameraDriver := must(camera.NewGstDriver(`
		autovideosrc ! video/x-raw,width=800,height=600,framerate=20/1 ! tee name=tee
			tee. ! queue ! valve name=valve-http-camera ! jpegenc ! appsink name=appsink-http-camera
			tee. ! queue ! valve name=valve-sip-camera ! videoscale ! videoconvert ! video/x-raw,format=I420,width=720,height=480 ! appsink name=appsink-sip-camera
	`))
	oplog.Default().SetWriter(&oplog.FileWriter{
		File: ".doorpix.oplog.jsonl",
	})

	var userAgent *sip.UserAgent = nil
	if cfg.SIP.Enabled {
		props := sip.UserAgentProps{
			Username:  cfg.SIP.Username,
			Password:  cfg.SIP.Password,
			Realm:     cfg.SIP.Realm,
			Domain:    cfg.SIP.Server,
			Webcam:    must(camera.NewWebcam("sip-camera", cameraDriver)),
			Whitelist: cfg.SIP.Whitelist,
		}

		slog.Debug("server: create sip user agent", "username", props.Username, "realm", props.Realm, "domain", props.Domain, "whitelist", props.Whitelist)
		userAgent = sip.NewUserAgent(props)
	}

	var httpServer *http.Server = nil
	if cfg.HTTP.Enabled {
		props := http.ServerProps{
			Webcam:    must(camera.NewWebcam("http-camera", cameraDriver)),
			UserAgent: userAgent,
			Port:      &cfg.HTTP.Port,
		}

		slog.Debug("server: create http server", "port", cfg.HTTP.Port)
		httpServer = http.NewServer(props)
	}

	shellService := system.NewShellService()

	return &Server{
		cameraDriver: cameraDriver,
		userAgent:    userAgent,
		httpServer:   httpServer,
		shellService: shellService,
	}
}

func (s *Server) Exec() {
	// Startup
	slog.Info("app server: starting up doorpix")
	oplog.Dispatch("system/doorpix/lifecycle/booting", "lifecycle", "booting")

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

	if s.shellService != nil {
		go func() {
			s.shellService.Serve()
		}()
	}

	oplog.Dispatch("system/doorpix/lifecycle/booted", "lifecycle", "booted")

	// Catch Signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Shutdown
	slog.Info("app server: shutting down doorpix")
	oplog.Dispatch("system/doorpix/lifecycle/stopping", "lifecycle", "stopping")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

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

	if s.shellService != nil {
		s.shellService.Shutdown(ctx)
	}

	oplog.Dispatch("system/doorpix/lifecycle/shutdown", "lifecycle", "shutdown")
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
