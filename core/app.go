package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/go-gst/go-gst/gst"
)

type App struct {
	system   doorpix.System
	handlers []Handler
}

func NewAppWithConfig(system doorpix.System) *App {
	app := &App{
		system: system,
	}

	app.init()

	return app
}

func (app *App) init() {
}

func (app *App) RegisterHandler(handler Handler) {
	app.handlers = append(app.handlers, handler)
}

func (app *App) setup() {
	for _, handler := range app.handlers {
		handler.Setup()
	}

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
		camera.NewElement("videoscale"),
		camera.NewElement(
			"capsfilter",
			"caps", gst.NewCapsFromString("video/x-raw,width=1920,height=1080"),
		),
		camera.NewElement("videoconvert"),
		camera.NewElement(
			"capsfilter",
			"caps", gst.NewCapsFromString("video/x-raw,format=I420,framerate=30/1"),
		),
	)

	if err != nil {
		panic(err)
	}

	return c
}
