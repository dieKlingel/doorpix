package doorpi

import (
	"context"
	"log/slog"

	"github.com/dieklingel/doorpix/internal/oplog"
)

type CallService struct {
	done      chan struct{}
	userAgent UserAgent
	channel   <-chan oplog.Event
}

func NewSipService(userAgent UserAgent) *CallService {
	return &CallService{
		done:      make(chan struct{}),
		userAgent: userAgent,
	}
}

func (s *CallService) Listen() {
	s.channel = oplog.On("internal/doorpix/service/call/invite")
}

func (s *CallService) Serve() error {
	for {
		select {
		case <-s.done:
			return nil
		case ev := <-s.channel:
			slog.Debug("call invite: received new invite event", "event", ev)
			event := &CallEvent{}

			err := oplog.UnmarshalEvent(ev.Properties, event)
			if err != nil {
				slog.Error("call invite: cannot process event", "error", err.Error())
				continue
			}

			_, err = s.userAgent.Invite(event.Uri)
			if err != nil {
				slog.Error("call invite: an error occoured inviting for a call", "error", err.Error())
			}
		}
	}
}

func (s *CallService) Run() error {
	if s.channel == nil {
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
