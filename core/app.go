package core

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type App struct {
	config doorpix.Config
	bus    *Bus

	services []Service
	wg       sync.WaitGroup
}

func NewAppWithConfig(config doorpix.Config, bus *Bus) *App {
	app := &App{
		config: config,
		bus:    bus,

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

	// exec all exec services
	ctx, cancel := context.WithCancel(ctx)
	app.exec(ctx)
	app.bus.Write(doorpix.StartupEvent, nil)

	// catch system signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	app.bus.Write(doorpix.ShutdownEvent, nil)
	app.bus.Close()
	cancel()
	app.wg.Wait()

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

	serviceNames := make([]string, len(app.services))
	for i, service := range app.services {
		serviceNames[i] = service.Name()
	}

	slog.Debug("initialized services", "services", serviceNames)
}

func (app *App) exec(ctx context.Context) {

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		for {
			event, ok := <-app.bus.Listen()
			if !ok {
				break
			}

			actions, ok := app.config.OnEvents[event.Type]
			if !ok {
				continue
			}

			hook := doorpix.NewActionHook(event.Data)
			for _, action := range actions {
				sucess := false
				for _, service := range app.services {
					if service, ok := service.(RunnerService); ok {
						sucess = service.Run(action, hook)
						if sucess {
							break
						}
					}
				}
				if !sucess {
					slog.Warn("no service could run the action", "action", action)
				}
			}
		}
	}()

	for _, service := range app.services {
		if service, ok := service.(ExecService); ok {
			if err := service.Exec(ctx, &app.wg); err != nil {
				panic(err)
			}
		}
	}
}
