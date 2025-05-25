package eventemitter

import "slices"

type Listener struct {
	Pattern string

	channel      chan Context
	eventEmitter *EventEmitter
}

func NewListener(pattern string) *Listener {
	return &Listener{
		Pattern: pattern,
		channel: make(chan Context),
	}
}

func (l *Listener) Listen() chan Context {
	return l.channel
}

func (l *Listener) Close() {
	index := slices.Index(l.eventEmitter.listeners, l)
	if index >= 0 {
		l.eventEmitter.listeners = slices.Delete(l.eventEmitter.listeners, index, index+1)
	}

	close(l.channel)
	l.eventEmitter = nil
}
