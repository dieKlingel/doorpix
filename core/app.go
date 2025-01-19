package core

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dieklingel/doorpix/core/internal/config"
	"github.com/dieklingel/doorpix/core/internal/exec"
)

type App struct {
	emitter *EventEmitter
}

func NewApp() *App {
	return &App{
		emitter: NewEventEmitter(),
	}
}

func (app *App) setup() {
	app.emitter.Listen(func(action config.Action, event *Event) {
		switch action := action.(type) {
		case config.LogAction:
			slog.Info(action.Message)
		case config.SleepAction:
			time.Sleep(time.Duration(action.Duration) * time.Second)
		case config.EvalAction:
			out, err := exec.Run(action.Expressions)
			if err != nil {
				slog.Error(err.Error())
				break
			}
			fmt.Print(out)
		}
	})

	app.emitter.Before(config.StartupEvent)
	app.emitter.On(config.StartupEvent)
	app.emitter.After(config.StartupEvent)
}

func (app *App) cleanup() {
	// cleanup the application state
	app.emitter.Before(config.ShutdownEvent)
	app.emitter.On(config.ShutdownEvent)
	app.emitter.After(config.ShutdownEvent)

	app.emitter.Wait()
}

func (app *App) Exec() {
	app.setup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	app.cleanup()
	os.Exit(1)
}
