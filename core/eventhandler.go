package core

import "github.com/dieklingel/doorpix/core/internal/config"

type EventHandler interface {
	HandleEvent(action config.Action, event *Event)
}
