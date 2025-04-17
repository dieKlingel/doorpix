package actions_test

import (
	"testing"
	"text/template"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/stretchr/testify/assert"
)

func TestConditionActionExecute(t *testing.T) {

	t.Run("plain false returns else", func(t *testing.T) {
		action := &actions.ConditionAction{
			Condition: *template.Must(template.New("condition").Parse("false")),
			Then:      nil,
			Else:      make([]actions.Action, 0),
		}

		result, err := action.Execute(nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("plain true returns then", func(t *testing.T) {
		action := &actions.ConditionAction{
			Condition: *template.Must(template.New("condition").Parse("true")),
			Then:      make([]actions.Action, 0),
			Else:      nil,
		}

		result, err := action.Execute(nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("evaluate to true returns then", func(t *testing.T) {
		action := &actions.ConditionAction{
			Condition: *template.Must(template.New("condition").Parse("[[ \"{{ .key }}\" == \"Hello\" ]]")),
			Then:      make([]actions.Action, 0),
			Else:      nil,
		}

		data := map[string]any{
			"key": "Hello",
		}
		result, err := action.Execute(data)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
