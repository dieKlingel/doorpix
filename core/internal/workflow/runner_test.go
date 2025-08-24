package workflow_test

import (
	"fmt"
	"path"
	"testing"

	"github.com/dieklingel/doorpix/core/internal/workflow"
	"github.com/dieklingel/doorpix/core/internal/workflow/event"
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

func TestRunnerFindPipelines(t *testing.T) {
	t.Run("should return all pipelines", func(t *testing.T) {
		runner := workflow.NewRunner(workflow.NewRegistry())
		runner.SetPipelineSource([]workflow.Pipeline{
			{
				Trigger: "test",
			},
		})

		pipelines, err := runner.FindPipelines(event.New("test"))
		assert.NoError(t, err)
		assert.Len(t, pipelines, 1)
	})

	t.Run("should return matching pipelines", func(t *testing.T) {
		runner := workflow.NewRunner(workflow.NewRegistry())
		runner.SetPipelineSource([]workflow.Pipeline{
			{
				Trigger: "test/*",
			},
			{
				Trigger: "test/*",
			},
			{
				Trigger: "other/*",
			},
		})

		pipelines, err := runner.FindPipelines(event.New("test/1"))
		assert.NoError(t, err)
		assert.Len(t, pipelines, 2)
	})

	t.Run("should handle erro in pipeline source", func(t *testing.T) {
		sourceErr := fmt.Errorf("source error")

		runner := workflow.NewRunner(workflow.NewRegistry())
		runner.SetPipelineSourceFunc(func() ([]workflow.Pipeline, error) {
			return nil, sourceErr
		})

		_, err := runner.FindPipelines(event.New("test/1"))
		assert.ErrorIs(t, err, sourceErr)
	})

	t.Run("should return err on unset source", func(t *testing.T) {
		runner := workflow.NewRunner(workflow.NewRegistry())
		_, err := runner.FindPipelines(event.New("test/1"))

		assert.ErrorIs(t, err, workflow.ErrSourceIsNil)
	})

	t.Run("should return err on invalid pattern", func(t *testing.T) {
		runner := workflow.NewRunner(workflow.NewRegistry())
		runner.SetPipelineSource([]workflow.Pipeline{
			{
				Trigger: "a[",
			},
		})

		_, err := runner.FindPipelines(event.New("test/1"))
		assert.ErrorIs(t, err, path.ErrBadPattern)
	})
}
