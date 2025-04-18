package actions

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v3"
)

type InviteAction struct {
	UriTemplates []template.Template
}

type Caller interface {
	Invite(uris []string) error
}

func (a *InviteAction) Execute(caller Caller, data map[string]any) error {
	uris := make([]string, len(a.UriTemplates))

	for i, uriTemplate := range a.UriTemplates {
		buffer := bytes.Buffer{}

		err := uriTemplate.Execute(&buffer, data)
		if err != nil {
			return err
		}

		uris[i] = buffer.String()
	}

	return caller.Invite(uris)
}

func (a *InviteAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		UriTemplates YamlScalarOrList[YamlStringTemplate] `yaml:"invite"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	a.UriTemplates = make([]template.Template, len(action.UriTemplates))
	for i := range action.UriTemplates {
		a.UriTemplates[i] = (template.Template)(action.UriTemplates[i])
	}
	return nil
}
