package server

import (
	"context"
	"errors"
	"reflect"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/dieklingel/doorpix/internal/oplog"
)

type EventMuxer struct {
	done   chan struct{}
	events []config.EventHandler
}

func NewEventMuxer(events []config.EventHandler) *EventMuxer {
	m := &EventMuxer{
		done:   make(chan struct{}),
		events: events,
	}

	return m
}

func (m *EventMuxer) Run() error {
	channels := map[int][]config.Step{}
	cases := make([]reflect.SelectCase, len(m.events)+1)

	for i, event := range m.events {
		ch := oplog.On(event.Event)
		steps := event.Steps

		channels[i] = steps

		c := reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
		cases[i] = c
	}

	cases[len(cases)-1] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(m.done),
	}

	remaining := len(cases)
	for remaining > 0 {
		chose, value, ok := reflect.Select(cases)
		if !ok {
			// The chosen channel has been closed, so zero out the channel to disable the case
			cases[chose].Chan = reflect.ValueOf(nil)
			remaining -= 1
			continue
		}

		steps, ok := channels[chose]
		// value received, but steps do not exists, so this was the done channel
		if !ok {
			return nil
		}

		in := value.Interface().(oplog.Event)
		for _, step := range steps {
			topic, exists := config.ServiceType[step.Type]
			if !exists {
				continue
			}

			args := make([]any, 0, len(step.Properties)*2+4)
			args = append(args, "parentId", in.Id, "source", "config")
			for key, value := range step.Properties {
				args = append(args, key, value)
			}

			oplog.Dispatch(topic, args...)
		}
	}

	return errors.New("eventmuxer: all channels closed")
}

func (m *EventMuxer) Stop(ctx context.Context) error {
	select {
	case m.done <- struct{}{}:
	case <-ctx.Done():
	}

	return nil
}
