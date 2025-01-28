package doorpix

type Bus interface {
	On(EventType)
	OnWithData(EventType, map[string]interface{})
	Handler(BusEventHandler)
	Wait()
}

type BusEventHandler interface {
	HandleEvent(Action, *ActionHook)
}
