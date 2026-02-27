package eventemitter

import (
	"fmt"
	"time"
)

type InMemoryEventEmitter struct {
	size  int
	index int
	log   []Event
}

func NewEventEmitter() EventEmitter {
	size := 30
	index := 0

	return &InMemoryEventEmitter{
		size:  size,
		index: index,
		log:   make([]Event, 0, size),
	}
}

func (e *InMemoryEventEmitter) Events() []Event {
	size := len(e.log)

	events := make([]Event, 0, size)
	events = append(events, e.log[e.index:size]...)
	events = append(events, e.log[0:e.index]...)

	return events
}

func (e *InMemoryEventEmitter) DispatchProperties(path string, properties map[string]any) (Event, error) {
	event := Event{
		Path:       path,
		Properties: properties,
		Timestamp:  time.Now(),
	}

	if len(e.log) < e.size {
		e.log = append(e.log, event)
	} else {
		e.log[e.index] = event
	}
	e.index++
	if e.index >= e.size {
		e.index = 0
	}

	return event, nil
}

func (e *InMemoryEventEmitter) Dispatch(path string, args ...any) (Event, error) {
	if len(args)%2 != 0 {
		return Event{}, ErrInvalidNumberOfArguments
	}

	values := make(map[string]any, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		arg := args[i]
		value := args[i+1]

		key, ok := arg.(string)
		if !ok {
			return Event{}, fmt.Errorf("%w: the key %s at position %d is not of type string", ErrInvalidArgumentType, key, i)
		}

		values[key] = value
	}

	return e.DispatchProperties(path, values)
}
