package config

import "gopkg.in/yaml.v3"

type Camera struct {
	Device string
}

func (camera *Camera) UnmarshalYAML(node *yaml.Node) error {
	raw := struct {
		Device *string
	}{}

	if err := node.Decode(&raw); err != nil {
		return err
	}

	if raw.Device != nil {
		camera.Device = *raw.Device
	}
	return nil
}
