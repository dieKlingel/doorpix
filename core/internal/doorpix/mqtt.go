package doorpix

import "gopkg.in/yaml.v3"

type MQTTConfig struct {
	Enabled       bool
	Host          string
	Port          int
	Protocol      string
	Username      string
	Password      string
	Subscribtions []string
}

func (mqttConfig *MQTTConfig) UnmarshalYAML(node *yaml.Node) error {
	rawHTTPConfig := struct {
		Enabled       bool     `yaml:"enabled"`
		Host          string   `yaml:"host"`
		Port          int      `yaml:"port"`
		Protocol      string   `yaml:"protocol"`
		Username      string   `yaml:"username"`
		Password      string   `yaml:"password"`
		Subscribtions []string `yaml:"subscribtions"`
	}{
		Enabled:  true,
		Port:     1883,
		Protocol: "tcp",
	}

	if err := node.Decode(&rawHTTPConfig); err != nil {
		return err
	}

	mqttConfig.Enabled = rawHTTPConfig.Enabled
	mqttConfig.Host = rawHTTPConfig.Host
	mqttConfig.Port = rawHTTPConfig.Port
	mqttConfig.Protocol = rawHTTPConfig.Protocol
	mqttConfig.Username = rawHTTPConfig.Username
	mqttConfig.Password = rawHTTPConfig.Password
	mqttConfig.Subscribtions = rawHTTPConfig.Subscribtions
	return nil
}
