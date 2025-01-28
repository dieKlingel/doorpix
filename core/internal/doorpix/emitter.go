package doorpix

type Emit = func(event EventType, data map[string]any)
