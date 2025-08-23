package yaml

import (
	"text/template"

	"gopkg.in/yaml.v3"
)

type Template template.Template

func (t *Template) UnmarshalYAML(node *yaml.Node) error {
	var val string

	err := node.Decode(&val)
	if err != nil {
		return err
	}

	tmpl, err := template.New("").Parse(val)
	if err != nil {
		return err
	}

	*t = Template(*tmpl)
	return nil
}

func (t *Template) Template() *template.Template {
	return (*template.Template)(t)
}
