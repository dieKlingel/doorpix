package doorpix

import "gopkg.in/yaml.v3"

type MQTTConfig struct {
	Enabled       bool
	ClientId      string
	Host          string
	Port          int
	Protocol      string
	Username      string
	Password      string
	Subscriptions []string
}

func (mqttConfig *MQTTConfig) UnmarshalYAML(node *yaml.Node) error {
	rawMQTTConfig := struct {
		Enabled       bool     `yaml:"enabled"`
		ClientId      string   `yaml:"clientId"`
		Host          string   `yaml:"host"`
		Port          int      `yaml:"port"`
		Protocol      string   `yaml:"protocol"`
		Username      string   `yaml:"username"`
		Password      string   `yaml:"password"`
		Subscriptions []string `yaml:"subscriptions"`
	}{
		Enabled:  true,
		Port:     1883,
		Protocol: "tcp",
	}

	if err := node.Decode(&rawMQTTConfig); err != nil {
		return err
	}

	mqttConfig.Enabled = rawMQTTConfig.Enabled
	mqttConfig.ClientId = rawMQTTConfig.ClientId
	mqttConfig.Host = rawMQTTConfig.Host
	mqttConfig.Port = rawMQTTConfig.Port
	mqttConfig.Protocol = rawMQTTConfig.Protocol
	mqttConfig.Username = rawMQTTConfig.Username
	mqttConfig.Password = rawMQTTConfig.Password
	mqttConfig.Subscriptions = rawMQTTConfig.Subscriptions
	return nil
}
