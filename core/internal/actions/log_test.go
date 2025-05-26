package actions_test

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/dieklingel/doorpix/core/internal/logs"
	"github.com/stretchr/testify/assert"
)

func TestLogActionExecute(t *testing.T) {

	t.Run("template without variabled", func(t *testing.T) {
		message := "test log message"
		action := actions.LogAction{
			Message: template.Must(template.New("").Parse(message)),
		}

		var buf bytes.Buffer
		err := action.Execute(&buf, nil)
		assert.NoError(t, err)

		assert.Equal(t, fmt.Sprintln(message), buf.String())
	})

	t.Run("template with variable", func(t *testing.T) {
		message := "test log message {{ .test }}"
		action := actions.LogAction{
			Message: template.Must(template.New("").Parse(message)),
		}

		var buf bytes.Buffer
		err := action.Execute(&buf, map[string]any{"test": "test"})
		assert.NoError(t, err)

		assert.Equal(t, fmt.Sprintln("test log message test"), buf.String())
	})

	t.Run("template with variable and nil map", func(t *testing.T) {
		message := "test log message {{ .test }}"
		action := actions.LogAction{
			Message: template.Must(template.New("").Parse(message)),
		}

		var buf bytes.Buffer
		err := action.Execute(&buf, nil)
		assert.NoError(t, err)

		assert.Equal(t, fmt.Sprintln("test log message "), buf.String())
	})

	t.Run("log to stdout", func(t *testing.T) {
		message := "test log message"
		action := actions.LogAction{
			Message: template.Must(template.New("").Parse(message)),
		}

		var output string
		logger := logs.IoWriterFunc(func(msg string) {
			output += msg
		})
		err := action.Execute(logger, nil)
		assert.Equal(t, fmt.Sprintln(message), output)
		assert.NoError(t, err)
	})
}
