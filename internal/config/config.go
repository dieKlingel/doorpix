package config

import (
	"bytes"
	"html/template"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HTTP   HTTP           `yaml:"http"`
	SIP    SIP            `yaml:"sip"`
	GPIO   GPIO           `yaml:"gpio"`
	Camera Camera         `yaml:"camera"`
	Events []EventHandler `yaml:"events"`
}

func Parse(content []byte) (*Config, error) {
	var c Config = Config{
		HTTP: HTTP{
			Enabled: true,
			Port:    8080,
		},
		GPIO: GPIO{
			Enabled:      false,
			DebounceTime: 200 * time.Millisecond,
		},
		Camera: Camera{
			Device: "autovideosrc",
		},
	}
	var expanded = bytes.Buffer{}

	err := template.Must(template.New("config").Parse(string(content))).Execute(&expanded, EnvVars())
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(expanded.Bytes(), &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func EnvVars() map[string]string {
	env := map[string]string{}

	for _, e := range os.Environ() {
		pair := bytes.SplitN([]byte(e), []byte("="), 2)
		if len(pair) != 2 {
			continue
		}

		env[string(pair[0])] = string(pair[1])
	}

	return env
}
