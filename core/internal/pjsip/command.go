package pjsip

import "sync"

type Command interface {
	Error() chan error
}

type BaseCommand struct {
	result chan error
	once   sync.Once
}

func (cmd *BaseCommand) Error() chan error {
	cmd.once.Do(func() {
		cmd.result = make(chan error)
	})

	return cmd.result
}

type SendInstantMessageCommand struct {
	BaseCommand

	Uri     string
	Message string
}
