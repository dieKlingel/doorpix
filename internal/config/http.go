package config

import "gopkg.in/yaml.v3"

type HTTP struct {
	Enabled bool
	Port    int
}

func (http *HTTP) UnmarshalYAML(node *yaml.Node) error {
	rawHTTPConfig := struct {
		Enabled bool `yaml:"enabled"`
		Port    int  `yaml:"port"`
	}{
		Enabled: true,
		Port:    8080,
	}

	if err := node.Decode(&rawHTTPConfig); err != nil {
		return err
	}

	http.Enabled = rawHTTPConfig.Enabled
	http.Port = rawHTTPConfig.Port
	return nil
}
