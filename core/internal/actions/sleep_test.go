package actions_test

import (
	"testing"
	"time"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/stretchr/testify/assert"
)

func TestSleepActionExecute(t *testing.T) {

	t.Run("test", func(t *testing.T) {
		sleeperCalled := false
		duration := 2 * time.Second

		sleepFunc := func(d time.Duration) {
			sleeperCalled = true
			assert.Equal(t, duration, d)
		}

		action := &actions.SleepAction{
			Duration: duration,
		}

		err := action.Execute(sleepFunc)
		assert.NoError(t, err)
		assert.True(t, sleeperCalled)
	})
}
