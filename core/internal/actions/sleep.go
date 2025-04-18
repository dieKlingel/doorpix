package actions

import (
	"time"

	"gopkg.in/yaml.v3"
)

type SleeperFunc func(time.Duration)

type SleepAction struct {
	Duration time.Duration `yaml:"sleep"`
}

func (a *SleepAction) Execute(sleeper SleeperFunc) error {
	sleeper(a.Duration)
	return nil
}

func (a *SleepAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		Sleep string `yaml:"sleep"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(action.Sleep)
	if err != nil {
		return err
	}

	a.Duration = duration
	return nil
}
