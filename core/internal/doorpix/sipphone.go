package doorpix

type SIPPhone struct {
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	Server      string   `yaml:"server"`
	Realm       string   `yaml:"realm"`
	StunServers []string `yaml:"stun-servers"`
	Whitelist   []string `yaml:"whitelist"`
}
