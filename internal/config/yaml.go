package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Bool bool

func (b *Bool) UnmarshalYAML(node *yaml.Node) error {
	var raw any
	if err := node.Decode(&raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case bool:
		*b = Bool(v)
	case string:
		switch v {
		case "true", "True", "TRUE", "":
			*b = Bool(true)
		case "false", "False", "FALSE":
			*b = Bool(false)
		default:
			return &yaml.TypeError{Errors: []string{"invalid boolean value: " + v}}
		}
	default:
		return &yaml.TypeError{Errors: []string{"expected boolean or string, got " + node.ShortTag()}}
	}

	return nil
}

type Int int

func (i *Int) UnmarshalYAML(node *yaml.Node) error {
	var raw any
	if err := node.Decode(&raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case int:
		*i = Int(v)
	case string:
		var parsed int
		if _, err := fmt.Sscanf(v, "%d", &parsed); err != nil {
			return &yaml.TypeError{Errors: []string{"invalid integer value: " + v}}
		}
		*i = Int(parsed)
	default:
		return &yaml.TypeError{Errors: []string{"expected integer or string, got " + node.ShortTag()}}
	}

	return nil
}
