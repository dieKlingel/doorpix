package core

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type App struct {
	system doorpix.System

	services []Service
	wg       sync.WaitGroup
}

func NewAppWithConfig(system doorpix.System) *App {
	app := &App{
		system: system,

		services: make([]Service, 0),
	}

	return app
}

func (app *App) RegisterService(service Service) {
	app.services = append(app.services, service)
}

func (app *App) Exec(ctx context.Context) {
	// init all init services
	app.init()

	app.system.Bus.On(doorpix.StartupEvent)

	// exec all exec services
	ctx, cancel := context.WithCancel(ctx)
	app.exec(ctx)

	// catch system signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	cancel()
	app.deinit()
	os.Exit(1)
}

func (app *App) init() {
	for _, service := range app.services {
		if service, ok := service.(InitService); ok {
			if err := service.Init(); err != nil {
				panic(err)
			}
		}
	}
}

func (app *App) exec(ctx context.Context) {
	for _, service := range app.services {
		if service, ok := service.(ExecService); ok {
			if err := service.Exec(ctx, &app.wg); err != nil {
				panic(err)
			}
		}
	}
}

func (app *App) deinit() {
	app.system.Bus.On(doorpix.ShutdownEvent)
	app.system.Bus.Wait()

	app.wg.Wait()
}
