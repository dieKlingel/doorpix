package eventemitter

type Listener struct {
	Pattern string
	Channel chan Context
}

func NewListener(pattern string) *Listener {
	return &Listener{
		Pattern: pattern,
		Channel: make(chan Context),
	}
}

func (l *Listener) Close() {
	close(l.Channel)
}
