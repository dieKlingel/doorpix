package oplog

type FilterFunc[T any] func(msg T) bool

type Listner[T any] struct {
	filters []FilterFunc[T]
	channel chan T
	Buffer  int
}

func (l *Listner[T]) Push(msg any) bool {
	event, ok := msg.(T)
	if !ok {
		return false
	}

	for _, filter := range l.filters {
		if !filter(event) {
			return false
		}
	}

	if l.channel == nil {
		l.channel = make(chan T, l.Buffer)
	}

	l.channel <- event
	return true
}

func (l *Listner[T]) FilterFunc(filter FilterFunc[T]) *Listner[T] {
	l.filters = append(l.filters, filter)
	return l
}

func (l *Listner[T]) Listen() <-chan T {
	if l.channel == nil {
		l.channel = make(chan T, l.Buffer)
	}

	return l.channel
}
