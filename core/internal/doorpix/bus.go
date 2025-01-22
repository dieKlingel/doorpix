package doorpix

type Bus interface {
	On(EventType)
	Handler(BusEventHandler)
	Wait()
}

type BusEventHandler interface {
	HandleEvent(Action, *Event)
}
