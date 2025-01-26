package doorpix

import "gopkg.in/yaml.v3"

type HTTPConfig struct {
	Enabled bool
	Port    int
}

func (httpConfig *HTTPConfig) UnmarshalYAML(node *yaml.Node) error {
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

	httpConfig.Enabled = rawHTTPConfig.Enabled
	httpConfig.Port = rawHTTPConfig.Port
	return nil
}
