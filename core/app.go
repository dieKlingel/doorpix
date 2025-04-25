package core

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/service"
)

type Softphone interface {
	actions.Caller
	actions.Messanger
}

type App struct {
	Config       doorpix.Config
	MQTTClient   actions.Publisher
	Softphone    Softphone
	EventEmitter *eventemitter.EventEmitter

	ctx service.Context
}

func (app *App) Start(ctx context.Context) error {
	slog.Info("Starting application")
	app.ctx = service.NewContext(context.Background())

	app.exec()
	return nil
}

func (app *App) Stop(ctx context.Context) error {
	slog.Info("Stopping application")

	app.ctx.CancelAndWait()

	return nil
}

func (app *App) exec() {
	listener := app.EventEmitter.Listen("events/*")
	app.ctx.Lock()
	go func() {
		defer app.ctx.Unlock()

		for {
			select {
			case <-app.ctx.Done():
				slog.Info("application context done")
				return
			case event := <-listener:
				eventType := strings.TrimPrefix(event.Event, "events/")
				workflow := app.Config.FindAllActionsByEventType(eventType)
				if len(workflow) == 0 {
					continue
				}

				app.executeWorklow(workflow, event.Data)
			}
		}
	}()
}

func (app *App) executeWorklow(workflow []actions.Action, ctx map[string]any) {
	for _, action := range workflow {
		var err error

		switch action := action.(type) {
		case actions.LogAction:
			err = action.Execute(os.Stdout, ctx)

		case actions.ConditionAction:
			actions, err := action.Execute(nil)
			if err == nil {
				app.executeWorklow(actions, ctx)
			}

		case actions.PublishAction:
			err = action.Execute(app.MQTTClient, ctx)

		case actions.MessageAction:
			err = action.Execute(app.Softphone, ctx)

		default:
			err = fmt.Errorf("unknown action type %T", action)
		}

		if err != nil {
			slog.Warn("failed to execute action", "action", action, "error", err)
		}
	}
}
