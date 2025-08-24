package oplog

type VerifyFunc func(msg any) error

type Logger struct {
	listeners []*Listner[any]
	verifiers []VerifyFunc
}

func NewLogger() *Logger {
	return &Logger{
		listeners: make([]*Listner[any], 0),
	}
}

func (l *Logger) Push(msg any) error {
	for _, verifier := range l.verifiers {
		if err := verifier(msg); err != nil {
			return err
		}
	}

	for _, listener := range l.listeners {
		listener.Push(msg)
	}

	return nil
}

func (l *Logger) AddVerifyFunc(verifier VerifyFunc) *Logger {
	l.verifiers = append(l.verifiers, verifier)
	return l
}

func (l *Logger) AddListener(listener *Listner[any]) *Logger {
	l.listeners = append(l.listeners, listener)
	return l
}
