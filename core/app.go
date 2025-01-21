package core

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dieklingel/doorpix/core/internal/config"
)

type App struct {
	emitter  *EventEmitter
	handlers []Handler
}

func NewApp() *App {
	emitter := NewEventEmitter()

	return &App{
		emitter: emitter,
	}
}

func (app *App) RegisterHandler(handler Handler) {
	app.handlers = append(app.handlers, handler)
}

func (app *App) setup() {
	for _, handler := range app.handlers {
		handler.Setup(app.emitter)
	}

	app.emitter.Before(config.StartupEvent)
	app.emitter.On(config.StartupEvent)
	app.emitter.After(config.StartupEvent)
}

func (app *App) cleanup() {
	// cleanup the application state
	app.emitter.Before(config.ShutdownEvent)

	for _, handler := range app.handlers {
		handler.Cleanup()
	}

	app.emitter.On(config.ShutdownEvent)
	app.emitter.After(config.ShutdownEvent)

	app.emitter.Wait()
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
