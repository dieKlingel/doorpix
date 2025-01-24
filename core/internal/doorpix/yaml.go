package doorpix

import (
	"text/template"

	"gopkg.in/yaml.v3"
)

type YamlScalarOrList[T any] []T

func (l *YamlScalarOrList[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		singleValueList := make(YamlScalarOrList[T], 1)
		if err := node.Decode(&singleValueList[0]); err != nil {
			return err
		}

		*l = singleValueList
		return nil
	}

	var list []T
	if err := node.Decode(&list); err != nil {
		return err
	}

	*l = list
	return nil
}

type YamlStringTemplate template.Template

func (t *YamlStringTemplate) UnmarshalYAML(node *yaml.Node) error {
	var val string

	err := node.Decode(&val)
	if err != nil {
		return err
	}

	tmpl, err := template.New("template").Parse(val)
	if err != nil {
		return err
	}

	*t = YamlStringTemplate(*tmpl)
	return nil
}
