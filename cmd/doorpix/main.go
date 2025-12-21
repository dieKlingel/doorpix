package main

import (
	"log/slog"

	"github.com/dieklingel/doorpix/internal/media/camera"
	"github.com/dieklingel/doorpix/internal/transport/http"
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

	httpServer.Serve()
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}
