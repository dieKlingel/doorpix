package config

type SIPPhone struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Realm    string `yaml:"realm"`
}
