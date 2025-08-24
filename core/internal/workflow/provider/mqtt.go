package provider

import (
	"bytes"

	"github.com/dieklingel/doorpix/core/internal/workflow"
	"github.com/dieklingel/doorpix/core/internal/yaml"
)

var MqttProviderTags = []string{
	"doorpix/mqtt.publish@v1",
}

type MqttClient interface {
	Publish(topic string, message string) error
}

type MqttProvider struct {
	client MqttClient
}

type MqttStepOptions struct {
	Topic   yaml.Template `yaml:"topic"`
	Message yaml.Template `yaml:"message"`
}

func NewMqttProvider(client MqttClient) *MqttProvider {
	return &MqttProvider{
		client: client,
	}
}

func (p *MqttProvider) Parse(step workflow.Step) (workflow.StepDelegate, error) {
	options := MqttStepOptions{}
	delegate := workflow.StepDelegate{
		Step:   step,
		Parent: options,
	}

	err := step.With.Decode(&options)
	if err != nil {
		return delegate, err
	}

	delegate.Execute = func(ctx *workflow.Context) error {
		topic := bytes.Buffer{}
		err := options.Topic.Template().Execute(&topic, ctx)
		if err != nil {
			return err
		}

		message := bytes.Buffer{}
		err = options.Message.Template().Execute(&message, ctx)
		if err != nil {
			return err
		}

		return p.client.Publish(topic.String(), message.String())
	}

	return delegate, nil
}
