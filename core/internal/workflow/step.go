package workflow

import (
	"gopkg.in/yaml.v3"
)

type Step struct {
	Uses string    `yaml:"uses"`
	With yaml.Node `yaml:"with"`
}

func NewStepWith(uses string, values map[string]any) (Step, error) {
	step := Step{
		Uses: uses,
		With: yaml.Node{},
	}

	err := step.With.Encode(values)
	if err != nil {
		return step, err
	}

	return step, nil
}
