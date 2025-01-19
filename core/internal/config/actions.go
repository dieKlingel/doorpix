package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Action interface{}

type SleepAction struct {
	Duration int `yaml:"sleep"`
}

type LogAction struct {
	Message string `yaml:"log"`
}

type EvalAction struct {
	Expressions []string `yaml:"eval"`
}

type InviteAction struct {
	Numbers []string `yaml:"invite"`
}

type HangupAction struct{}

func newActionFromNode(node yaml.Node) (Action, error) {
	if node.Kind == yaml.MappingNode {
		raw := map[string]any{}
		if err := node.Decode(&raw); err != nil {
			return nil, err
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
	} else if node.Kind == yaml.ScalarNode {
		if node.Value == "hangup" {
			return HangupAction{}, nil
		}
	}

	return nil, fmt.Errorf("could not infer action type: %s", node.Value)
}
