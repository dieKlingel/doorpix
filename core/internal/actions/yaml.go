package actions

import (
	"fmt"
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

func newActionFromNode(node yaml.Node) (Action, error) {
	if node.Kind == yaml.MappingNode {
		raw := map[string]any{}
		if err := node.Decode(&raw); err != nil {
			return nil, err
		}

		if raw["condition"] != nil {
			action := ConditionAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["sleep"] != nil {
			action := SleepAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["log"] != nil {
			action := LogAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["hangup"] != nil {
			action := HangupAction{}
			err := node.Decode(&action)
			return HangupAction{}, err
		}
		if raw["invite"] != nil {
			action := InviteAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["eval"] != nil {
			action := EvalAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["message"] != nil && raw["to"] != nil {
			action := MessageAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["topic"] != nil && raw["message"] != nil {
			action := PublishAction{}
			err := node.Decode(&action)
			return action, err
		}
		/*if raw["rpc"] != nil {
			action := RPCAction{}
			err := node.Decode(&action)
			return action, err
		}*/
	} else if node.Kind == yaml.ScalarNode {
		if node.Value == "hangup" {
			return HangupAction{}, nil
		}
	}

	return nil, fmt.Errorf("could not infer action type in line %d", node.Line)
}
