package eventemitter

type EventEmitter interface {
	Dispatch(path string, args ...any) (Event, error)
	DispatchProperties(path string, properties map[string]any) (Event, error)
	Events() []Event
	On(path string) chan Event
}
