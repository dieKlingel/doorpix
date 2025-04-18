package actions

type HangupAction struct{}

type Call interface {
	Hangup() error
}

func (a *HangupAction) Execute(call Call, _ map[string]any) error {
	return call.Hangup()
}
