package doorpi

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dieklingel/doorpix/internal/oplog"
)

type ShellService struct {
	commander Commander
	done      chan struct{}
}

func NewShellService(commander Commander) *ShellService {
	return &ShellService{
		commander: commander,
		done:      make(chan struct{}),
	}
}

func (s *ShellService) Run() error {
	channel := oplog.On("internal/doorpix/service/shell")

	for {
		select {
		case <-s.done:
			return nil
		case ev := <-channel:
			slog.Debug("system shell: received new shell event", "event", ev)
			event := &ShellEvent{
				Silent: true,
			}

			err := oplog.UnmarshalEvent(ev.Properties, event)
			if err != nil {
				slog.Error("system shell: cannot process event", "error", err.Error())
				continue
			}

			go func() {
				output, err := s.commander.Exec("sh", "-c", event.Cmd)
				if err != nil {
					slog.Error("system shell: an error occoured executing a command", "error", err.Error(), "command", event.Cmd, "silent", event.Silent)
				}
				if !event.Silent {
					fmt.Print(string(output))
				}
			}()
		}
	}
}

func (s *ShellService) Stop(ctx context.Context) error {
	select {
	case s.done <- struct{}{}:
	case <-ctx.Done():
	}

	return nil
}
