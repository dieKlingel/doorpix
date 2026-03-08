package eventemitter

import (
	"errors"
	"fmt"
	"math/big"
	"path"
	"time"
)

type InMemoryEventEmitter struct {
	revision  *big.Int
	size      int
	index     int
	log       []Event
	listeners map[string][]chan Event
}

func NewEventEmitter() EventEmitter {
	size := 30
	index := 0

	return &InMemoryEventEmitter{
		revision:  big.NewInt(0),
		size:      size,
		index:     index,
		log:       make([]Event, 0, size),
		listeners: make(map[string][]chan Event),
	}
}

func (e *InMemoryEventEmitter) Events() []Event {
	size := len(e.log)

	events := make([]Event, 0, size)
	events = append(events, e.log[e.index:size]...)
	events = append(events, e.log[0:e.index]...)

	return events
}

func (e *InMemoryEventEmitter) DispatchEvent(event Event) error {
	errs := make([]error, 0)

	for k, v := range e.listeners {
		match, err := path.Match(k, event.Path)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if match {
			for _, channel := range v {
				channel <- event
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}

func (e *InMemoryEventEmitter) DispatchProperties(path string, properties map[string]any) (Event, error) {
	rev := e.revision.Add(e.revision, big.NewInt(1))
	id := big.NewInt(1)
	id.Set(rev)

	event := Event{
		Id:         id.Text(10),
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

	e.DispatchEvent(event)
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

func (e *InMemoryEventEmitter) On(path string) <-chan Event {
	channel := make(chan Event)
	if _, exists := e.listeners[path]; !exists {
		e.listeners[path] = make([]chan Event, 0)
	}

	e.listeners[path] = append(e.listeners[path], channel)
	return channel
}
