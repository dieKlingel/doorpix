package doorpix_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfigUnmarshalFromYaml(t *testing.T) {
	t.Run("actions is parsed as yaml", func(t *testing.T) {
		data := `
on:
  startup:
    - sleep: 1s
`
		c := doorpix.NewConfig()
		err := yaml.Unmarshal([]byte(data), &c)

		assert.NoError(t, err)
		assert.Len(t, c.OnEvents[doorpix.StartupEvent], 1)
	})

	t.Run("action is parsed as type", func(t *testing.T) {
		data := `
on:
  startup:
    - sleep: 1s
    - log: hello
    - hangup: {}
`
		c := doorpix.NewConfig()
		err := yaml.Unmarshal([]byte(data), &c)
		assert.NoError(t, err)

		sleepAction := c.OnEvents[doorpix.StartupEvent][0]
		assert.IsType(t, actions.SleepAction{}, sleepAction)
		logAction := c.OnEvents[doorpix.StartupEvent][1]
		assert.IsType(t, actions.LogAction{}, logAction)
		hangupAction := c.OnEvents[doorpix.StartupEvent][2]
		assert.IsType(t, actions.HangupAction{}, hangupAction)
	})
}
