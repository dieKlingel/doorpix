package main

import (
	"log/slog"
	"os"

	"github.com/dieklingel/doorpix/core/internal/config"
)

func main() {
	config.AddConfigPath(
		"/etc/doorpix/doorpix.yaml",
		"/etc/doorpix/config.yaml",
		"doorpix.yaml",
		"config.yaml",
	)
	if err := config.Read(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
