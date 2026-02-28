package config

var ServiceType = map[string]string{
	"shell": "internal/doorpix/service/shell",
}

type Step struct {
	Type       string         `yaml:"type"`
	Properties map[string]any `yaml:"with"`
}

type EventHandler struct {
	Event string `yaml:"event"`
	Steps []Step `yaml:"steps"`
}
