package system

import (
	"github.com/dieklingel/doorpix/internal/oplog"
)

type ShellEvent struct {
	Cmd    string
	Silent bool
}

func NewShellEventFromEvent(event oplog.Event) (*ShellEvent, error) {
	cmd, err := oplog.ParseString(event, "cmd")
	if err != nil {
		return nil, err
	}
	silent := oplog.ParseBoolOrDefault(event, "silent", true)

	return &ShellEvent{
		Cmd:    cmd,
		Silent: silent,
	}, nil
}
