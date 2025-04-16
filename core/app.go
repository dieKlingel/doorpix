package core

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dieklingel/doorpix/core/internal/actionrunner"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/service/httpsvc"
	"github.com/dieklingel/doorpix/core/internal/service/mqttsvc"
)

type App struct {
	Config doorpix.Config

	registry   *actionrunner.Registry
	wg         sync.WaitGroup
	eventQueue *EventQueue
}

func NewApp(config doorpix.Config) *App {
	registry := actionrunner.Registry{}
	eventQueue := NewEventQueue()

	return &App{
		Config:     config,
		registry:   &registry,
		eventQueue: eventQueue,
	}
}

func (app *App) Exec(ctx context.Context) {
	// init all init services
	services := make([]Service, 0)

	if app.Config.HTTP.Enabled {
		http := httpsvc.New(httpsvc.HTTPServiceProps{
			Port:                    app.Config.HTTP.Port,
			VideoStreamCameraDevice: app.Config.Camera.Device,
		})
		services = append(services, http)
	}

	if app.Config.MQTT.Enabled {
		mqtt := mqttsvc.New(
			mqttsvc.MQTTServiceProps{
				Host:          app.Config.MQTT.Host,
				Port:          app.Config.MQTT.Port,
				Protocol:      app.Config.MQTT.Protocol,
				Username:      app.Config.MQTT.Username,
				Password:      app.Config.MQTT.Password,
				Subscriptions: app.Config.MQTT.Subscriptions,
			},
		)
		services = append(services, mqtt)
		app.registry.Publish = mqtt.Publish
	}

	for _, service := range services {
		if service, ok := service.(InitService); ok {
			err := service.Init()
			if err != nil {
				panic(err)
			}
		}
	}

	// exec all exec services
	ctx, cancel := context.WithCancel(ctx)

	app.exec(ctx)
	for _, service := range services {
		if service, ok := service.(BackgroundService); ok {
			err := service.StartBackgroundTask(ctx, &app.wg)
			if err != nil {
				panic(err)
			}
		}
	}

	app.eventQueue.Write(doorpix.StartupEvent, nil)

	// catch system signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	app.eventQueue.Close()
	cancel()
	app.wg.Wait()

	os.Exit(1)
}

func (app *App) exec(ctx context.Context) {

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()

		for {
			event, ok := <-app.eventQueue.Listen()
			if !ok {
				break
			}

			actions := app.Config.FindAllActionsByEventType(event.Type)
			if len(actions) == 0 {
				continue
			}

			for _, action := range actions {
				runnable, err := app.registry.CreateRunnable(action)
				if err != nil {
					slog.Error("could not create runnable", "error", err)
					continue
				}

				err = runnable.Run(ctx)
				if err != nil {
					slog.Error("could not run action", "error", err)
					continue
				}
			}
		}
	}()
}
