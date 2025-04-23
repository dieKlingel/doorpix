package providers

import (
	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

func NewApplicationConfiguration() doorpix.Config {
	config := doorpix.NewConfig()
	config.AddConfigPath(
		"/etc/doorpix/doorpix.yaml",
		"/etc/doorpix/config.yaml",
		"doorpix.yaml",
		"config.yaml",
	)
	if err := config.Read(); err != nil {
		panic(err)
	}
	if err := config.Error(); err != nil {
		panic(err)
	}

	return config
}
