package env

import (
	"fmt"
	"log/slog"
	"os"
)

func LogLevel() slog.Level {
	level := os.Getenv("DOORPIX_LOG_LEVEL")
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}

	if len(level) != 0 {
		slog.Warn(fmt.Sprintf("'%s' is not a valid log level", level))
	}

	return slog.LevelInfo
}
