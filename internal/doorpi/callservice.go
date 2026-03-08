package doorpi

import (
	"context"
	"log/slog"

	"github.com/dieklingel/doorpix/internal/oplog"
)

type CallService struct {
	done      chan struct{}
	userAgent UserAgent
}

func NewSipService(userAgent UserAgent) *CallService {
	return &CallService{
		done:      make(chan struct{}),
		userAgent: userAgent,
	}
}

func (s *CallService) Run() error {
	channel := oplog.On("internal/doorpix/service/call/invite")
	for {
		select {
		case <-s.done:
			return nil
		case ev := <-channel:
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

func (s *CallService) Stop(ctx context.Context) error {
	select {
	case s.done <- struct{}{}:
	case <-ctx.Done():
	}

	return nil
}
