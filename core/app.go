package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
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
