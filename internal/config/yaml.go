package config

import (
	"errors"
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
	case int:
		switch v {
		case 0:
			*b = Bool(false)
		case 1:
			*b = Bool(true)
		default:
			return errors.Join(
				ErrUnexpectedInput,
				&yaml.TypeError{Errors: []string{fmt.Sprintf("invalid int value: %d", v)}},
			)
		}
	case string:
		switch v {
		case "true", "True", "TRUE":
			*b = Bool(true)
		case "false", "False", "FALSE":
			*b = Bool(false)
		default:
			return errors.Join(
				ErrUnexpectedInput,
				&yaml.TypeError{Errors: []string{fmt.Sprintf("invalid boolean value: %s", v)}},
			)
		}
	default:
		return errors.Join(
			ErrUnexpectedInput,
			&yaml.TypeError{Errors: []string{fmt.Sprintf("expected boolean, int or string, got %s", node.ShortTag())}},
		)
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
		return errors.Join(
			ErrUnexpectedInput,
			&yaml.TypeError{Errors: []string{"expected integer or string, got " + node.ShortTag()}},
		)
	}

	return nil
}
