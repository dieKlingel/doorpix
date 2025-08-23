package workflow

import (
	"gopkg.in/yaml.v3"
)

type Step struct {
	Uses string    `yaml:"uses"`
	With yaml.Node `yaml:"with"`
}
