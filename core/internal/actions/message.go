package actions

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v3"
)

type MessageAction struct {
	UriTemplates    []template.Template `yaml:"to"`
	MessageTemplate template.Template   `yaml:"message"`
}

type Messanger interface {
	SendMessage(uris []string, message string) error
}

func (a *MessageAction) Execute(messanger Messanger, data map[string]any) error {
	uris := make([]string, len(a.UriTemplates))

	for i, uri := range a.UriTemplates {
		buffer := bytes.Buffer{}
		err := uri.Execute(&buffer, data)
		if err != nil {
			return err
		}
		uris[i] = buffer.String()
	}

	message := bytes.Buffer{}
	err := a.MessageTemplate.Execute(&message, data)
	if err != nil {
		return err
	}

	return messanger.SendMessage(uris, message.String())
}

func (a *MessageAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		UriTemplates YamlScalarOrList[YamlStringTemplate] `yaml:"to"`
		Message      YamlStringTemplate                   `yaml:"message"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	a.UriTemplates = make([]template.Template, len(action.UriTemplates))
	for i, number := range action.UriTemplates {
		a.UriTemplates[i] = (template.Template)(number)
	}
	a.MessageTemplate = (template.Template)(action.Message)
	return nil
}
