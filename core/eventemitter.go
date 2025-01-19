package core

import (
	"sync"

	"github.com/dieklingel/doorpix/core/internal/config"
)

type EventCallback = func(config.Action, *Event)

type EventEmitter struct {
	conf      config.Config
	listeners []EventCallback
	mutex     sync.Mutex
	waitgroup sync.WaitGroup
}

func NewEventEmitter() *EventEmitter {
	return NewEventEmitterWithConfig(*config.GetGlobal())
}

func NewEventEmitterWithConfig(conf config.Config) *EventEmitter {
	return &EventEmitter{
		conf:      conf,
		listeners: make([]EventCallback, 0),
	}
}

func (emitter *EventEmitter) Listen(cb EventCallback) {
	emitter.mutex.Lock()
	defer emitter.mutex.Unlock()

	emitter.listeners = append(emitter.listeners, cb)
}

func (emitter *EventEmitter) Execute(eventtype config.Event, actions []config.Action) {
	emitter.waitgroup.Add(1)
	go func() {
		defer emitter.waitgroup.Done()

		event := &Event{}
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

func (emitter *EventEmitter) Before(eventtype config.Event) {
	actions := emitter.conf.BeforeEvents[eventtype]
	if actions == nil {
		return
	}

	emitter.Execute(eventtype, actions)
}

func (emitter *EventEmitter) On(eventtype config.Event) {
	actions := emitter.conf.OnEvents[eventtype]
	if actions == nil {
		return
	}

	emitter.Execute(eventtype, actions)
}

func (emitter *EventEmitter) After(eventtype config.Event) {
	actions := emitter.conf.AfterEvents[eventtype]
	if actions == nil {
		return
	}

	emitter.Execute(eventtype, actions)
}
