package actions_test

import (
	"testing"
	"text/template"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/stretchr/testify/assert"
)

func TestEvalActionExecute(t *testing.T) {
	t.Run("execute non existing command", func(t *testing.T) {
		action := &actions.EvalAction{
			Expressions: []*template.Template{
				template.Must(template.New("").Parse("command_which_does_not_exist")),
			},
		}
		err := action.Execute(map[string]any{})
		assert.Error(t, err)
	})

	t.Run("execute command with exit code 1", func(t *testing.T) {
		action := &actions.EvalAction{
			Expressions: []*template.Template{
				template.Must(template.New("").Parse("exit 1")),
			},
		}
		err := action.Execute(map[string]any{})
		assert.Error(t, err)
	})

	t.Run("execute command with exit code 0", func(t *testing.T) {

		action := &actions.EvalAction{
			Expressions: []*template.Template{
				template.Must(template.New("").Parse("exit 0")),
			},
		}
		err := action.Execute(map[string]any{})
		assert.NoError(t, err)
	})
}
