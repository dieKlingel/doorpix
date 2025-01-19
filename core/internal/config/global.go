package config

var global = New()

func AddConfigPath(file ...string) {
	global.AddConfigPath(file...)
}

func Read() error {
	return global.Read()
}

func GetGlobal() *Config {
	return global
}
