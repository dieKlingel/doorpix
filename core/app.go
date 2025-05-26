package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/service"
)

type Softphone interface {
	actions.Caller
	actions.Messanger
	actions.Call
}

type App struct {
	Config              doorpix.Config
	EventEmitter        *eventemitter.EventEmitter
	BasePipelineBuilder *PipelineBuilder

	ctx service.Context
}

func (app *App) Start() error {
	slog.Info("Starting application")
	app.ctx = service.NewContext(context.Background())

	return app.exec()
}

func (app *App) Stop() {
	slog.Info("Stopping application")
	app.ctx.CancelAndWait()
}

func (app *App) exec() error {
	for eventPath, workflow := range app.Config.Workflows() {
		eventPath = fmt.Sprintf("events/%s", eventPath)
		listener := app.EventEmitter.Listen(eventPath)
		pipeline, err := app.BasePipelineBuilder.Clone().WithActions(workflow).Build()
		if err != nil {
			return fmt.Errorf("failed to build pipeline for event %s: %w", eventPath, err)
		}
		if pipeline.IsEmpty() {
			slog.Debug("workflow is empty, skipping", "eventPath", eventPath)
			continue
		}

		app.ctx.Lock()
		go func() {
			defer app.ctx.Unlock()

			for {
				select {
				case <-app.ctx.Done():
					return
				case event := <-listener.Listen():
					pipeline.Run(event.Data)

					// publish state updates e.g. used by the kiosk.
					// publishState at the end of each workflow for now, may be changed later
					app.EventEmitter.Emit("state/update", event.Data)
				}
			}
		}()
	}

	return nil
}
