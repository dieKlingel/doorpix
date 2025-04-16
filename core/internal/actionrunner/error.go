package actionrunner

import (
	"fmt"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type CouldNotCreateRunnableError struct {
	action doorpix.Action
}

func newCouldNotCreateRunnableError(action doorpix.Action) error {
	return &CouldNotCreateRunnableError{
		action: action,
	}
}

func (e *CouldNotCreateRunnableError) Error() string {
	return fmt.Sprintf("could not create runnable for action %T: %v", e.action, e.action)
}

type UnknownActionError struct {
	action doorpix.Action
}

func newUnknownActionError(action doorpix.Action) error {
	return &UnknownActionError{
		action: action,
	}
}

func (e *UnknownActionError) Error() string {
	return fmt.Sprintf("unknown action %T: %v", e.action, e.action)
}
