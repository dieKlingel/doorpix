package core

import (
	"sync"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type EventCallback = func(doorpix.Action, *doorpix.Event)

type EventEmitter struct {
	config    *doorpix.Config
	listeners []EventCallback
	mutex     sync.Mutex
	waitgroup sync.WaitGroup
}

func NewEventEmitterWithConfig(conf *doorpix.Config) *EventEmitter {
	return &EventEmitter{
		config:    conf,
		listeners: make([]EventCallback, 0),
	}
}

func (emitter *EventEmitter) Listen(cb EventCallback) {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()

	emitter.listeners = append(emitter.listeners, cb)
}

func (emitter *EventEmitter) Handler(handler doorpix.BusEventHandler) {
	emitter.Listen(handler.HandleEvent)
}

func (emitter *EventEmitter) Execute(eventtype doorpix.EventType, actions []doorpix.Action, data map[string]any) {
	emitter.waitgroup.Add(1)
	go func() {
		defer emitter.waitgroup.Done()

		event := doorpix.NewEvent()
		event.AddData(data)

		for _, action := range actions {
			for _, handler := range emitter.listeners {
				handler(action, event)
			}
		}
	}()
}

func (emitter *EventEmitter) Wait() {
	emitter.waitgroup.Wait()
}

func (emitter *EventEmitter) Before(eventtype doorpix.EventType) {
	actions := emitter.config.BeforeEvents[eventtype]
	if actions == nil {
		return
	}

	emitter.Execute(eventtype, actions, nil)
}

func (emitter *EventEmitter) On(eventtype doorpix.EventType) {
	actions := emitter.config.OnEvents[eventtype]
	if actions == nil {
		return
	}

	emitter.Execute(eventtype, actions, nil)
}

func (emitter *EventEmitter) OnWithData(eventtype doorpix.EventType, data map[string]any) {
	actions := emitter.config.OnEvents[eventtype]
	if actions == nil {
		return
	}

	emitter.Execute(eventtype, actions, data)
}

func (emitter *EventEmitter) After(eventtype doorpix.EventType) {
	actions := emitter.config.AfterEvents[eventtype]
	if actions == nil {
		return
	}

	emitter.Execute(eventtype, actions, nil)
}
