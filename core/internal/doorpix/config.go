package doorpix

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	files []string

	OnEvents     EventCollection `yaml:"on"`
	BeforeEvents EventCollection `yaml:"before"`
	AfterEvents  EventCollection `yaml:"after"`

	SIPPhone SIPPhone     `yaml:"sip-phone"`
	HTTP     HTTPConfig   `yaml:"http"`
	Camera   CameraConfig `yaml:"camera"`
	RPC      RPCConfig    `yaml:"rpc"`
	MQTT     MQTTConfig   `yaml:"mqtt"`
}

func NewConfig() Config {
	return Config{
		OnEvents:     EventCollection{},
		BeforeEvents: EventCollection{},
		AfterEvents:  EventCollection{},

		SIPPhone: SIPPhone{
			Enabled: false,
		},
		HTTP: HTTPConfig{
			Enabled: false,
		},
	}
}

func (c *Config) Error() error {
	if err := c.Camera.Error(); err != nil {
		return err
	}

	return nil
}

func (c *Config) FindAllActionsByEventType(eventtype EventType) []Action {
	return c.OnEvents[eventtype]
}

func (c *Config) AddConfigPath(file ...string) {
	c.files = append(c.files, file...)
}

func (c *Config) Read() error {
	for _, file := range c.files {
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				continue
			} else {
				return err
			}
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(content, c); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("no config file found in %v", c.files)
}
