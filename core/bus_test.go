package core_test

import (
	"sync"
	"testing"

	"github.com/dieklingel/doorpix/core"
	"github.com/stretchr/testify/assert"
)

func TestBus(t *testing.T) {
	t.Run("should write event", func(t *testing.T) {
		bus := core.NewEventQueue()

		channel := bus.Listen()
		go func() {
			event, ok := <-channel
			assert.True(t, ok)
			assert.Equal(t, "test", event.Type)
		}()

		bus.Write("test", nil)
	})

	t.Run("should close bus", func(t *testing.T) {
		bus := core.NewEventQueue()

		wg := sync.WaitGroup{}

		channel := bus.Listen()
		wg.Add(1)
		go func() {
			defer wg.Done()

			event := <-channel
			assert.Equal(t, "test", event.Type)

			_, ok := <-channel
			assert.False(t, ok)
		}()

		bus.Write("test", nil)
		bus.Close()
		wg.Wait()
	})
}
