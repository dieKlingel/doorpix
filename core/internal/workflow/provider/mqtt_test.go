package provider_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/workflow"
	"github.com/dieklingel/doorpix/core/internal/workflow/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMqttClient struct {
	mock.Mock
}

func (c *MockMqttClient) Publish(topic string, message string) error {
	args := c.Called(topic, message)
	return args.Error(0)
}

func TestMqttProviderParse(t *testing.T) {
	t.Run("should parse step with valid options", func(t *testing.T) {
		mqttClient := &MockMqttClient{}
		provider := provider.NewMqttProvider(mqttClient)

		options := map[string]any{
			"topic":   "test/topic",
			"message": "test message",
		}

		step, _ := workflow.NewStepWith("test", options)
		_, err := provider.Parse(step)
		assert.NoError(t, err)
	})

	t.Run("should execute succesfully", func(t *testing.T) {
		mqttClient := &MockMqttClient{}
		provider := provider.NewMqttProvider(mqttClient)

		options := map[string]any{
			"topic":   "test/topic",
			"message": "test message",
		}

		step, _ := workflow.NewStepWith("test", options)
		delegate, _ := provider.Parse(step)

		mqttClient.On("Publish", "test/topic", "test message").Return(nil)
		err := delegate.Execute(workflow.NewContext())
		assert.NoError(t, err)
	})
}
