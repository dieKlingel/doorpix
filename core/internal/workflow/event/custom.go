package event

type CustomEvent struct {
	value string
}

func New(value string) *CustomEvent {
	return &CustomEvent{
		value: value,
	}
}

func (e *CustomEvent) String() string {
	return e.value
}
