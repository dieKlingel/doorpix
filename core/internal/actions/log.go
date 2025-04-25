package actions

import (
	"html/template"
	"io"

	"gopkg.in/yaml.v3"
)

type LogAction struct {
	Message *template.Template
}

func (action *LogAction) Execute(writer io.Writer, data map[string]any) error {
	err := action.Message.Execute(writer, data)
	if err == nil {
		io.WriteString(writer, "\n")
	}
	return err
}

func (a *LogAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		Log string `yaml:"log"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	tmpl, err := template.New("log").Parse(action.Log)
	if err != nil {
		return err
	}

	a.Message = tmpl
	return nil
}
