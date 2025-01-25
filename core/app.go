package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/drivers/appvideo"
	"github.com/go-gst/go-gst/gst"
)

type App struct {
	system   doorpix.System
	handlers []Handler
}

func NewAppWithConfig(system doorpix.System) *App {
	return &App{
		system: system,
	}
}

func (app *App) RegisterHandler(handler Handler) {
	app.handlers = append(app.handlers, handler)
}

func (app *App) setup() {
	for _, handler := range app.handlers {
		handler.Setup()
	}
	appvideo.SetCameraFactory(app.newCameraFactory)
	appvideo.Initialize()

	app.system.Bus.On(doorpix.StartupEvent)
}

func (app *App) cleanup() {
	// cleanup the application state
	app.system.Bus.On(doorpix.ShutdownEvent)

	for _, handler := range app.handlers {
		handler.Cleanup()
	}

	app.system.Bus.Wait()
}

func (app *App) Exec() {
	app.setup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for _, handler := range app.handlers {
		handler.Exec()
	}

	<-c
	app.cleanup()
	os.Exit(1)
}

func (app *App) newCameraFactory() *camera.Camera {
	c, err := camera.NewFromString(
		app.system.Config.Camera.Device,
		camera.NewElement(
			"capsfilter",
			"caps", gst.NewCapsFromString("video/x-raw,format=I420,width=640,height=480,framerate=25/1"),
		),
	)

	if err != nil {
		panic(err)
	}

	return c
}
