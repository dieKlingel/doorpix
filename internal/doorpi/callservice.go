package doorpi

import (
	"context"
	"log/slog"

	"github.com/dieklingel/doorpix/internal/oplog"
)

type CallService struct {
	done      chan struct{}
	userAgent UserAgent
	invites   <-chan oplog.Event
	messages  <-chan oplog.Event
}

func NewSipService(userAgent UserAgent) *CallService {
	return &CallService{
		done:      make(chan struct{}),
		userAgent: userAgent,
	}
}

func (s *CallService) Listen() {
	s.invites = oplog.On("internal/doorpix/service/call/invite")
	s.messages = oplog.On("internal/doorpix/service/call/message")
}

func (s *CallService) Serve() error {
	for {
		select {
		case <-s.done:
			return nil
		case input := <-s.invites:
			slog.Debug("call invite: received new invite event", "event", input)
			event := &CallEvent{}

			err := oplog.UnmarshalEvent(input.Properties, event)
			if err != nil {
				slog.Error("call invite: cannot process event", "error", err.Error())
				continue
			}

			_, err = s.userAgent.Invite(event.Uri)
			if err != nil {
				slog.Error("call invite: an error occoured inviting for a call", "error", err.Error())
			}
		case input := <-s.messages:
			slog.Debug("call message: received new message event", "event", input)
			event := &MessageEvent{}

			err := oplog.UnmarshalEvent(input.Properties, event)
			if err != nil {
				slog.Error("call message: cannot process event", "error", err.Error())
				continue
			}

			err = s.userAgent.SendMessage(event.Uri, event.Body)
			if err != nil {
				slog.Error("call message: an error occoured sending a message", "error", err.Error(), "uri", event.Uri)
			}
		}
	}
}

func (s *CallService) Run() error {
	if s.invites == nil {
		s.Listen()
	}

	s.Serve()
	return nil
}

func (s *CallService) Stop(ctx context.Context) error {
	select {
	case s.done <- struct{}{}:
	case <-ctx.Done():
	}

	return nil
}
