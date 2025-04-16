package actionrunner

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type LogActionHandler func(message string) error
type EvalActionHandler func(expressions []string) error
type PublishActionHandler func(topic string, message string) error
type SleepActionHandler func(duration time.Duration) error
type SipMessageHandler func(uris []string, message string) error
type HangupMessageHandler func() error
type InviteMessageHandler func(uris []string) error

type Registry struct {
	Log        LogActionHandler
	Eval       EvalActionHandler
	Publish    PublishActionHandler
	Sleep      SleepActionHandler
	SipMessage SipMessageHandler
	Hangup     HangupMessageHandler
	Invite     InviteMessageHandler
}

func (registry *Registry) Builder() *ActionRunnerRegistryBuilder {
	return &ActionRunnerRegistryBuilder{
		registry: registry,
	}
}

type ActionRunnerRegistryBuilder struct {
	registry *Registry
}

type Runnable struct {
	Run func(ctx context.Context) error
}

func newRunnable(cb func(ctx context.Context) error) *Runnable {
	return &Runnable{
		Run: cb,
	}
}

func (builder *Registry) CreateRunnable(action doorpix.Action) (*Runnable, error) {
	switch action := action.(type) {
	case doorpix.LogAction:
		return builder.createLogRunnable(&action)
	case doorpix.EvalAction:
		return builder.createEvalRunnable(&action)
	case doorpix.PublishAction:
		return builder.createPublishRunnable(&action)
	case doorpix.SleepAction:
		return builder.createSleepRunnable(&action)
	case doorpix.MessageAction:
		return builder.createSipMessageRunnable(&action)
	case doorpix.HangupAction:
		return builder.createHangupRunnable(&action)
	case doorpix.InviteAction:
		return builder.createInviteRunnable(&action)
	case doorpix.ConditionAction:
		return builder.createConditionRunnable(&action)
	}

	return nil, newUnknownActionError(action)
}

func (registry *Registry) createLogRunnable(action *doorpix.LogAction) (*Runnable, error) {
	if registry.Log == nil {
		return nil, newCouldNotCreateRunnableError(action)
	}

	runnable := &Runnable{
		Run: func(ctx context.Context) error {
			buffer := bytes.Buffer{}

			err := action.Message.Execute(&buffer, ctx)
			if err != nil {
				return err
			}

			err = registry.Log(buffer.String())
			return err
		},
	}

	return runnable, nil
}

func (registry *Registry) createEvalRunnable(action *doorpix.EvalAction) (*Runnable, error) {
	if registry.Eval == nil {
		return nil, newCouldNotCreateRunnableError(action)
	}

	runnable := newRunnable(func(ctx context.Context) error {
		expressions := make([]string, len(action.Expressions))

		for i, expr := range action.Expressions {
			buffer := bytes.Buffer{}

			err := expr.Execute(&buffer, ctx)
			if err != nil {
				return err
			}

			expressions[i] = buffer.String()
		}

		err := registry.Eval(expressions)
		return err
	})

	return runnable, nil
}

func (registry *Registry) createPublishRunnable(action *doorpix.PublishAction) (*Runnable, error) {
	if registry.Publish == nil {
		return nil, newCouldNotCreateRunnableError(action)
	}

	runnable := newRunnable(func(ctx context.Context) error {
		buffer := bytes.Buffer{}

		err := action.Topic.Execute(&buffer, ctx)
		if err != nil {
			return err
		}

		topic := buffer.String()

		buffer.Reset()
		err = action.Message.Execute(&buffer, ctx)
		if err != nil {
			return err
		}

		message := buffer.String()

		err = registry.Publish(topic, message)
		return err
	})

	return runnable, nil
}

func (registry *Registry) createSleepRunnable(action *doorpix.SleepAction) (*Runnable, error) {
	if registry.Sleep == nil {
		return nil, newCouldNotCreateRunnableError(action)
	}

	runnable := newRunnable(func(ctx context.Context) error {
		return registry.Sleep(action.Duration)
	})

	return runnable, nil
}

func (registry *Registry) createSipMessageRunnable(action *doorpix.MessageAction) (*Runnable, error) {
	if registry.SipMessage == nil {
		return nil, newCouldNotCreateRunnableError(action)
	}

	runnable := newRunnable(func(ctx context.Context) error {
		buffer := bytes.Buffer{}

		uris := make([]string, len(action.UriTemplates))
		for i, uri := range action.UriTemplates {
			err := uri.Execute(&buffer, ctx)
			if err != nil {
				return err
			}
			uris[i] = buffer.String()
			buffer.Reset()
		}

		err := action.MessageTemplate.Execute(&buffer, ctx)
		if err != nil {
			return err
		}

		message := buffer.String()
		err = registry.SipMessage(uris, message)
		return err
	})

	return runnable, nil
}

func (registry *Registry) createHangupRunnable(action *doorpix.HangupAction) (*Runnable, error) {
	if registry.Hangup == nil {
		return nil, newCouldNotCreateRunnableError(action)
	}

	runnable := newRunnable(func(ctx context.Context) error {
		return registry.Hangup()
	})

	return runnable, nil
}

func (registry *Registry) createInviteRunnable(action *doorpix.InviteAction) (*Runnable, error) {
	if registry.Invite == nil {
		return nil, newCouldNotCreateRunnableError(action)
	}

	runnable := newRunnable(func(ctx context.Context) error {
		buffer := bytes.Buffer{}

		uris := make([]string, len(action.UriTemplates))
		for i, uri := range action.UriTemplates {
			err := uri.Execute(&buffer, ctx)
			if err != nil {
				return err
			}
			uris[i] = buffer.String()
			buffer.Reset()
		}

		err := registry.Invite(uris)
		return err
	})

	return runnable, nil
}

func (registry *Registry) createConditionRunnable(action *doorpix.ConditionAction) (*Runnable, error) {
	runnable := newRunnable(func(ctx context.Context) error {
		buffer := bytes.Buffer{}

		err := action.Condition.Execute(&buffer, ctx)
		if err != nil {
			return err
		}

		condition := buffer.String()
		conditionIsTrue, err := strconv.ParseBool(condition)
		if err != nil {
			return err
		}

		var furtherActions []doorpix.Action
		if conditionIsTrue {
			furtherActions = action.Then
		} else {
			furtherActions = action.Else
		}

		for _, furtherAction := range furtherActions {
			runnable, err := registry.CreateRunnable(furtherAction)
			if err != nil {
				return err
			}

			err = runnable.Run(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return runnable, nil
}
