package doorpix

import "gopkg.in/yaml.v3"

type SIPPhone struct {
	Enabled     bool     `yaml:"enabled"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	Server      string   `yaml:"server"`
	Realm       string   `yaml:"realm"`
	StunServers []string `yaml:"stun-servers"`
	Whitelist   []string `yaml:"whitelist"`
}

func (phone *SIPPhone) UnmarshalYAML(node *yaml.Node) error {
	raw := struct {
		Enabled     bool     `yaml:"enabled"`
		Username    string   `yaml:"username"`
		Password    string   `yaml:"password"`
		Server      string   `yaml:"server"`
		Realm       string   `yaml:"realm"`
		StunServers []string `yaml:"stun-servers"`
		Whitelist   []string `yaml:"whitelist"`
	}{
		Enabled: true,
	}

	if err := node.Decode(&raw); err != nil {
		return err
	}

	phone.Enabled = raw.Enabled
	phone.Username = raw.Username
	phone.Password = raw.Password
	phone.Server = raw.Server
	phone.Realm = raw.Realm
	phone.StunServers = raw.StunServers
	phone.Whitelist = raw.Whitelist
	return nil
}
