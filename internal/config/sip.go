package config

import "gopkg.in/yaml.v3"

type SIP struct {
	Enabled     bool
	Username    string
	Password    string
	Server      string
	Realm       string
	StunServers []string
	Whitelist   []string
}

func (sip *SIP) UnmarshalYAML(node *yaml.Node) error {
	raw := struct {
		Enabled     Bool     `yaml:"enabled"`
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

	sip.Enabled = bool(raw.Enabled)
	sip.Username = raw.Username
	sip.Password = raw.Password
	sip.Server = raw.Server
	sip.Realm = raw.Realm
	sip.StunServers = raw.StunServers
	sip.Whitelist = raw.Whitelist
	return nil
}
