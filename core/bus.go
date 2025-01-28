package core

import (
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type Event struct {
	Type doorpix.EventType
	Data map[string]interface{}
}

type Bus struct {
	listener chan Event
}

func NewBus() *Bus {
	return &Bus{
		listener: make(chan Event),
	}
}

func (bus *Bus) Listen() chan Event {
	return bus.listener
}

func (bus *Bus) Write(t doorpix.EventType, data map[string]interface{}) {
	slog.Debug("bus received event", "type", t, "data", data)

	event := Event{
		Type: t,
		Data: data,
	}

	if bus.listener != nil {
		bus.listener <- event
	}
	slog.Debug("bus distributed event", "type", t, "data", data)
}

func (bus *Bus) Close() {
	close(bus.listener)
}
