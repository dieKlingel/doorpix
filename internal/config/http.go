package config

import "gopkg.in/yaml.v3"

type HTTP struct {
	Enabled bool
	Port    int
}

func (http *HTTP) UnmarshalYAML(node *yaml.Node) error {
	raw := struct {
		Enabled *Bool `yaml:"enabled"`
		Port    *Int  `yaml:"port"`
	}{}

	if err := node.Decode(&raw); err != nil {
		return err
	}

	if raw.Enabled != nil {
		http.Enabled = bool(*raw.Enabled)
	}
	if raw.Port != nil {
		http.Port = int(*raw.Port)
	}
	return nil
}
