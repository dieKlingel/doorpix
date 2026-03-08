package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/dieklingel/doorpix/internal/device/gpio"
	"github.com/dieklingel/doorpix/internal/device/shell"
	"github.com/dieklingel/doorpix/internal/doorpi"
	"github.com/dieklingel/doorpix/internal/media/camera"
	"github.com/dieklingel/doorpix/internal/oplog"
	"github.com/dieklingel/doorpix/internal/transport/http"
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

type Server struct {
	workers []Worker
}

func New(cfg *config.Config) *Server {
	workers := make([]Worker, 0, 6)
	cameraDriver := must(camera.NewGstDriver(fmt.Sprintf(`
		%s ! tee name=tee
			tee. ! queue ! valve name=valve-http-camera ! jpegenc ! appsink name=appsink-http-camera
			tee. ! queue ! valve name=valve-sip-camera ! videoscale ! videoconvert ! video/x-raw,format=I420,width=720,height=480 ! appsink name=appsink-sip-camera
		`,
		cfg.Camera.Device,
	)))

	oplog.Default().SetWriter(&oplog.FileWriter{
		File: "logs/.doorpix.oplog.jsonl",
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
		workers = append(workers, userAgent)
	}

	if cfg.HTTP.Enabled {
		props := http.ServerProps{
			Webcam:    must(camera.NewWebcam("http-camera", cameraDriver)),
			UserAgent: userAgent,
			Port:      &cfg.HTTP.Port,
		}

		slog.Debug("server: create http server", "port", cfg.HTTP.Port)
		workers = append(workers, http.NewServer(props))
	}

	if cfg.GPIO.Enabled {
		props := gpio.ControllerProps{
			Chip:         cfg.GPIO.Chip,
			Inputs:       cfg.GPIO.Inputs,
			DebounceTime: cfg.GPIO.DebounceTime,
		}

		slog.Debug("server: create gpio controller", "chip", props.Chip, "debounce-time", props.DebounceTime, "inputs", props.Inputs)
		workers = append(workers, gpio.NewController(props))
	}

	if userAgent != nil {
		workers = append(workers, doorpi.NewSipService(userAgent))
	}
	workers = append(workers, doorpi.NewShellService(shell.NewController()))
	workers = append(workers, NewEventMuxer(cfg.Events))

	return &Server{
		workers: workers,
	}
}

func (s *Server) Exec() {
	// Startup
	slog.Info("app server: starting up doorpix")
	oplog.Dispatch("system/doorpix/lifecycle/booting", "lifecycle", "booting")

	for _, worker := range s.workers {
		go func() {
			err := worker.Run()
			if err != nil {
				slog.Error("an error occoured running a worker", "error", err.Error())
			}
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
	wg := sync.WaitGroup{}

	for _, worker := range s.workers {
		wg.Go(func() {
			err := worker.Stop(ctx)
			if err != nil {
				slog.Error("an error occoured stopping a worker", "error", err.Error())
			}
		})
	}

	wg.Wait()
	oplog.Dispatch("system/doorpix/lifecycle/shutdown", "lifecycle", "shutdown")
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
