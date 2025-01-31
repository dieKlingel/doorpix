package doorpix_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestLogAction_UnmarshalYaml(t *testing.T) {
	t.Run("should return error on invalid yaml", func(t *testing.T) {
		data := `log`

		var action doorpix.LogAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.Error(t, err)
	})

	t.Run("should parse string", func(t *testing.T) {
		data := `log: "test"`

		var action doorpix.LogAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.NoError(t, err)
		assert.NotNil(t, action.Message)
	})

	t.Run("should return error on invalid template", func(t *testing.T) {
		data := `log: "{{}"`

		var action doorpix.LogAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.Error(t, err)
	})

	t.Run("should parse template", func(t *testing.T) {
		data := `log: "{{ .Test }}"`

		var action doorpix.LogAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.NoError(t, err)

		var buf bytes.Buffer
		err = action.Message.Execute(&buf, map[string]interface{}{
			"Test": "test",
		})
		assert.NoError(t, err)
		assert.Equal(t, "test", buf.String())
	})

}

func TestEvalAction_UnmarshalYaml(t *testing.T) {
	t.Run("should return error on invalid yaml", func(t *testing.T) {
		data := `eval`

		var action doorpix.EvalAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.Error(t, err)
	})

	t.Run("should parse empty array", func(t *testing.T) {
		data := `eval: []`

		var action doorpix.EvalAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.NoError(t, err)
		assert.Len(t, action.Expressions, 0)
	})

	t.Run("should parse array with one element", func(t *testing.T) {
		data := `eval: ["1 + 2"]`

		var action doorpix.EvalAction
		err := yaml.Unmarshal([]byte(data), &action)

		assert.NoError(t, err)
		assert.Len(t, action.Expressions, 1)

		var bytes bytes.Buffer
		err = action.Expressions[0].Execute(&bytes, nil)
		assert.NoError(t, err)
		assert.Equal(t, "1 + 2", bytes.String())
	})

	t.Run("should parse string", func(t *testing.T) {
		data := `eval: "1 + 2"`

		var action doorpix.EvalAction
		err := yaml.Unmarshal([]byte(data), &action)

		assert.NoError(t, err)
		assert.Len(t, action.Expressions, 1)

		var bytes bytes.Buffer
		err = action.Expressions[0].Execute(&bytes, nil)
		assert.NoError(t, err)
		assert.Equal(t, "1 + 2", bytes.String())
	})
}

func TestSleepAction_UnmarshalYaml(t *testing.T) {
	t.Run("should parse time", func(t *testing.T) {
		data := `sleep: 1000ms`

		var action doorpix.SleepAction
		err := yaml.Unmarshal([]byte(data), &action)

		assert.NoError(t, err)
		assert.Equal(t, 1000*time.Millisecond, action.Duration)
	})

	t.Run("should return error on int", func(t *testing.T) {
		data := `sleep: 1000`

		var action doorpix.SleepAction
		err := yaml.Unmarshal([]byte(data), &action)

		assert.Error(t, err)
	})
}

func TestInviteAction_UnmarshalYaml(t *testing.T) {
	t.Run("should parse string", func(t *testing.T) {
		data := `invite: "sip:1234@localhost"`

		var action doorpix.InviteAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.NoError(t, err)
	})

	t.Run("should parse array", func(t *testing.T) {
		data := `invite: ["sip:1234@localhost", "sip:5678@localhost"]`

		var action doorpix.InviteAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.NoError(t, err)
	})

	t.Run("should parse template", func(t *testing.T) {
		data := `invite: "{{ .Test }}"`

		var action doorpix.InviteAction
		err := yaml.Unmarshal([]byte(data), &action)
		assert.NoError(t, err)

		var buf bytes.Buffer
		err = action.UriTemplates[0].Execute(&buf, map[string]interface{}{
			"Test": "sip:1234@localhost",
		})
		assert.NoError(t, err)
		assert.Equal(t, "sip:1234@localhost", buf.String())
	})
}
