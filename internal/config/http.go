package config

import "gopkg.in/yaml.v3"

type HTTP struct {
	Enabled bool
	Port    int
}

func (http *HTTP) UnmarshalYAML(node *yaml.Node) error {
	rawHTTPConfig := struct {
		Enabled Bool `yaml:"enabled"`
		Port    Int  `yaml:"port"`
	}{
		Enabled: Bool(true),
		Port:    Int(8080),
	}

	if err := node.Decode(&rawHTTPConfig); err != nil {
		return err
	}

	http.Enabled = bool(rawHTTPConfig.Enabled)
	http.Port = int(rawHTTPConfig.Port)
	return nil
}
