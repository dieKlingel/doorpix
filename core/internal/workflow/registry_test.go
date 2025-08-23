package workflow_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/workflow"
	"github.com/dieklingel/doorpix/core/internal/workflowtest"
	"github.com/stretchr/testify/assert"
)

func TestRegistryRegisterProvider(t *testing.T) {

	t.Run("should register a new provider successfully", func(t *testing.T) {

		registry := workflow.NewRegistry()
		provider := &workflowtest.NilProvider{}

		err := registry.RegisterProvider("doorpix/mock@v1", provider)
		assert.NoError(t, err)

		retrievedProvider, exists := registry.GetProvider("doorpix/mock@v1")
		assert.True(t, exists)
		assert.Equal(t, provider, retrievedProvider)
	})

	t.Run("should return an error when registering a provider with an existing identifier", func(t *testing.T) {
		registry := workflow.NewRegistry()
		provider := &workflowtest.NilProvider{}

		err := registry.RegisterProvider("doorpix/mock@v1", provider)
		assert.NoError(t, err)

		err = registry.RegisterProvider("doorpix/mock@v1", provider)
		assert.ErrorIs(t, err, workflow.ErrProviderAlreadyRegistered)
	})

	t.Run("should return false when retrieving a non-existent provider", func(t *testing.T) {
		registry := workflow.NewRegistry()

		_, exists := registry.GetProvider("doorpix/non-existent@v1")
		assert.False(t, exists)
	})

	t.Run("should allow register multiple tags for the same provider", func(t *testing.T) {
		registry := workflow.NewRegistry()
		provider := &workflowtest.NilProvider{}

		err := registry.RegisterProvider("doorpix/mock@v1", provider)
		assert.NoError(t, err)

		err = registry.RegisterProvider("doorpix/mock@v2", provider)
		assert.NoError(t, err)

		_, exists := registry.GetProvider("doorpix/mock@v1")
		assert.True(t, exists)

		_, exists = registry.GetProvider("doorpix/mock@v2")
		assert.True(t, exists)
	})
}
