package core_test

import (
	"testing"
	"text/template"
	"time"

	"github.com/dieklingel/doorpix/core"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

func TestEventLoopEmitBefore(t *testing.T) {
	t.Run("should emit before event", func(t *testing.T) {
		config := doorpix.NewConfig()
		config.BeforeEvents[doorpix.StartupEvent] = append(config.BeforeEvents[doorpix.StartupEvent], doorpix.LogAction{
			Message: template.Must(template.New("log").Parse("before event")),
		})

		callback := make(chan bool, 1)

		emitter := core.NewEventEmitterWithConfig(config)
		emitter.Listen(func(action doorpix.Action, event *doorpix.Event) {
			if _, ok := action.(doorpix.LogAction); ok {
				callback <- true
			}
		})

		emitter.Before(doorpix.StartupEvent)

		select {
		case <-callback:
		case <-time.After(1 * time.Second):
			t.Error("the callback was not called")
		}
	})
}
