package workflow_test

import (
	"fmt"
	"testing"

	"github.com/dieklingel/doorpix/core/internal/workflow"
	"github.com/dieklingel/doorpix/core/internal/workflowtest"
	"github.com/stretchr/testify/assert"
)

func TestRunnerRun(t *testing.T) {
	t.Run("should return error for missing registry", func(t *testing.T) {
		runner := workflow.NewRunner(nil)
		err := runner.Run(nil)

		assert.ErrorIs(t, err, workflow.ErrRegistryIsNil)
	})

	t.Run("should return error for nil pipeline", func(t *testing.T) {
		runner := workflow.NewRunner(workflow.NewRegistry())
		err := runner.Run(nil)

		assert.ErrorIs(t, err, workflow.ErrPipelineIsNil)
	})

	t.Run("should run successfull with empty pipeline", func(t *testing.T) {
		runner := workflow.NewRunner(workflow.NewRegistry())
		err := runner.Run(&workflow.Pipeline{})

		assert.NoError(t, err)
	})

	t.Run("should return error for missing provider", func(t *testing.T) {
		runner := workflow.NewRunner(workflow.NewRegistry())
		err := runner.Run(&workflow.Pipeline{
			Steps: []workflow.Step{
				{
					Uses: "doorpix/non-existent-provider",
				},
			},
		})

		assert.ErrorIs(t, err, workflow.ErrProviderNotFound)
	})

	t.Run("should return error for provider parse error", func(t *testing.T) {
		nilProvider := &workflowtest.NilProvider{
			ParseError: fmt.Errorf("parse error"),
		}
		registry := workflow.NewRegistry()
		registry.RegisterProvider("doorpix/nil-provider", nilProvider)

		runner := workflow.NewRunner(registry)
		err := runner.Run(&workflow.Pipeline{
			Steps: []workflow.Step{
				{
					Uses: "doorpix/nil-provider",
				},
			},
		})

		assert.Error(t, err)
	})

	t.Run("should return error for provider execution error", func(t *testing.T) {
		nilProvider := &workflowtest.NilProvider{
			RunError: fmt.Errorf("run error"),
		}
		registry := workflow.NewRegistry()
		registry.RegisterProvider("doorpix/nil-provider", nilProvider)

		runner := workflow.NewRunner(registry)
		err := runner.Run(&workflow.Pipeline{
			Steps: []workflow.Step{
				{
					Uses: "doorpix/nil-provider",
				},
			},
		})

		assert.Error(t, err)
	})
}
