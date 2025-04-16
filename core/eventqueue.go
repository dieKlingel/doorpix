package core

import (
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type Event struct {
	Type doorpix.EventType
	Data map[string]interface{}
}

type EventQueue struct {
	listener chan Event
}

func NewEventQueue() *EventQueue {
	return &EventQueue{
		listener: make(chan Event),
	}
}

func (queue *EventQueue) Listen() chan Event {
	return queue.listener
}

func (queue *EventQueue) Write(t doorpix.EventType, data map[string]interface{}) {
	slog.Debug("queue received event", "type", t, "data", data)

	event := Event{
		Type: t,
		Data: data,
	}

	if queue.listener != nil {
		queue.listener <- event
	}
	slog.Debug("queue distributed event", "type", t, "data", data)
}

func (bus *EventQueue) Close() {
	close(bus.listener)
}
