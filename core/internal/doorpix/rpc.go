package doorpix

import "gopkg.in/yaml.v3"

type RPCConfig struct {
	Enabled bool
	Host    string
	Port    int
}

func (rpc *RPCConfig) UnmarshalYAML(node *yaml.Node) error {
	raw := struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
	}{
		Enabled: true,
		Host:    "127.0.0.1",
		Port:    50051,
	}

	if err := node.Decode(&raw); err != nil {
		return err
	}

	rpc.Enabled = raw.Enabled
	rpc.Host = raw.Host
	rpc.Port = raw.Port

	return nil
}
