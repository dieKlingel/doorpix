package actions

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v3"
)

type PublishAction struct {
	Topic   template.Template
	Message template.Template
}

type Publisher interface {
	Publish(topic string, message string) error
}

func (a *PublishAction) Execute(publisher Publisher, data map[string]any) error {
	topic := bytes.Buffer{}
	err := a.Topic.Execute(&topic, data)
	if err != nil {
		return err
	}

	message := bytes.Buffer{}
	err = a.Message.Execute(&message, data)
	if err != nil {
		return err
	}

	return publisher.Publish(topic.String(), message.String())
}

func (a *PublishAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		Topic   YamlStringTemplate `yaml:"topic"`
		Message YamlStringTemplate `yaml:"message"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	a.Topic = (template.Template)(action.Topic)
	a.Message = (template.Template)(action.Message)
	return nil
}
