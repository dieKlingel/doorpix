package config

var ServiceType = map[string]string{
	"shell":   "internal/doorpix/service/shell",
	"invite":  "internal/doorpix/service/call/invite",
	"message": "internal/doorpix/service/call/message",
}

type Step struct {
	Type       string         `yaml:"type"`
	Properties map[string]any `yaml:"with"`
}

type EventHandler struct {
	Event string `yaml:"event"`
	Steps []Step `yaml:"steps"`
}
