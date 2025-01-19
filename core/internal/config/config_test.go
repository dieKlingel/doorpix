package config_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestConfigUnmarshalFromYaml(t *testing.T) {
	t.Run("actions is parsed as yaml", func(t *testing.T) {
		data := `
on:
  startup:
    - sleep: 1
`
		c := config.New()
		err := yaml.Unmarshal([]byte(data), &c)

		assert.NoError(t, err)
		assert.Len(t, c.OnEvents[config.StartupEvent], 1)
	})

	t.Run("action is parsed as type", func(t *testing.T) {
		data := `
on:
  startup:
    - sleep: 1
    - log: hello
    - hangup: {}
`
		c := config.New()
		err := yaml.Unmarshal([]byte(data), &c)
		assert.NoError(t, err)

		sleepAction := c.OnEvents[config.StartupEvent][0]
		assert.IsType(t, config.SleepAction{}, sleepAction)
		logAction := c.OnEvents[config.StartupEvent][1]
		assert.IsType(t, config.LogAction{}, logAction)
		hangupAction := c.OnEvents[config.StartupEvent][2]
		assert.IsType(t, config.HangupAction{}, hangupAction)
	})
}
