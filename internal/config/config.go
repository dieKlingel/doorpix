package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP HTTP `yaml:"http"`
	SIP  SIP  `yaml:"sip"`
}

func Parse(content []byte) (*Config, error) {
	var c Config
	expanded := os.ExpandEnv(string(content))
	println(expanded)

	err := yaml.Unmarshal([]byte(expanded), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
