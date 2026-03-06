package config

import (
	"time"

	"gopkg.in/yaml.v3"
)

type GPIO struct {
	Enabled      bool
	Chip         string
	DebounceTime time.Duration
	Inputs       []int
}

func (gpio *GPIO) UnmarshalYAML(node *yaml.Node) error {
	raw := struct {
		Enabled      Bool           `yaml:"enabled"`
		Chip         string         `yaml:"chip"`
		DebounceTime *time.Duration `yaml:"debounce-time"`
		Inputs       []Int          `yaml:"inputs"`
	}{
		Enabled: true,
	}

	if err := node.Decode(&raw); err != nil {
		return err
	}

	gpio.Enabled = bool(raw.Enabled)
	gpio.Chip = raw.Chip
	gpio.Inputs = make([]int, len(raw.Inputs))
	for i, input := range raw.Inputs {
		gpio.Inputs[i] = int(input)
	}
	if raw.DebounceTime != nil {
		gpio.DebounceTime = *raw.DebounceTime
	}

	return nil
}
