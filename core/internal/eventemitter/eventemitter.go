package eventemitter

import (
	"log/slog"
	"path"
)

type Context struct {
	Event string
	Data  map[string]any
}

type EventEmitter struct {
	listeners []*Listener
}

func New() *EventEmitter {
	return &EventEmitter{
		listeners: make([]*Listener, 0),
	}
}

func (e *EventEmitter) EmitOptional(event string, data map[string]any) error {
	for _, listener := range e.listeners {
		match, err := path.Match(listener.Pattern, event)
		if err != nil {
			return err
		}

		if match {
			listener.channel <- Context{
				Event: event,
				Data:  data,
			}
		}
	}

	return nil
}

func (e *EventEmitter) Emit(event string, data map[string]any) error {
	listenrCount := 0

	for _, listener := range e.listeners {
		match, err := path.Match(listener.Pattern, event)
		if err != nil {
			return err
		}

		if match {
			listener.channel <- Context{
				Event: event,
				Data:  data,
			}
			listenrCount++
		}

	}

	if listenrCount == 0 {
		slog.Warn("no listeners for event", "event", event)
	}

	return nil
}

func (e *EventEmitter) Listen(pattern string) *Listener {
	listener := NewListener(pattern)
	listener.eventEmitter = e
	e.listeners = append(e.listeners, listener)

	return listener
}
