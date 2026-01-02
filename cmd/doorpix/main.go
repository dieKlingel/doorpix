package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	_ "github.com/dieklingel/doorpix/pkg/pjsua2"

	"github.com/dieklingel/doorpix/internal/media/camera"
	"github.com/dieklingel/doorpix/internal/transport/http"
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("starting doorpix...")

	driver := must(camera.NewGstDriver(`
		autovideosrc ! video/x-raw,width=800,height=600,framerate=20/1 ! tee name=tee
			tee. ! queue ! valve name=valve-http-camera ! jpegenc ! appsink name=appsink-http-camera
	`))

	httpServer := http.NewServer(http.ServerProps{
		Webcam: must(camera.NewWebcam("http-camera", driver)),
	})
	sipClient := sip.NewClient(sip.ClientProps{})

	serve(&httpServer)
	serve(&sipClient)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	slog.Info("stopping doorpix")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	shutdown(&httpServer, ctx)
	shutdown(&sipClient, ctx)
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

type Server interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

func serve(srv Server) {
	go func() {
		err := srv.Serve()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
}

func shutdown(srv Server, ctx context.Context) {
	err := srv.Shutdown(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
}
