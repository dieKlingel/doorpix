package workflow_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/workflow"
	"github.com/stretchr/testify/assert"
)

func TestNewStepWith(t *testing.T) {
	t.Run("should create step with options", func(t *testing.T) {
		options := map[string]any{
			"key": "value",
		}

		step, err := workflow.NewStepWith("test", options)

		assert.NoError(t, err)
		assert.Equal(t, "test", step.Uses)
		assert.NotNil(t, step.With)
	})
}
