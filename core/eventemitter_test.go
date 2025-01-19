package core_test

import (
	"testing"
	"text/template"
	"time"

	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/config"
)

func TestEventLoopEmitBefore(t *testing.T) {
	t.Run("should emit before event", func(t *testing.T) {
		conf := config.New()
		conf.BeforeEvents[config.StartupEvent] = append(conf.BeforeEvents[config.StartupEvent], config.LogAction{
			Message: template.Must(template.New("log").Parse("before event")),
		})

		callback := make(chan bool, 1)

		emitter := core.NewEventEmitterWithConfig(conf)
		emitter.Listen(func(action config.Action, event *core.Event) {
			if _, ok := action.(config.LogAction); ok {
				callback <- true
			}
		})

		emitter.Before(config.StartupEvent)

		select {
		case <-callback:
		case <-time.After(1 * time.Second):
			t.Error("the callback was not called")
		}
	})
}
