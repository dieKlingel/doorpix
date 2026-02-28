package system

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/dieklingel/doorpix/internal/oplog"
)

type ShellService struct {
	done chan struct{}
}

func NewShellService() *ShellService {
	return &ShellService{
		done: make(chan struct{}),
	}
}

func (s *ShellService) Serve() {
	channel := oplog.On("internal/doorpix/service/shell")

	for {
		select {
		case <-s.done:
			return
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
				output, err := exec.Command("sh", "-c", event.Cmd).CombinedOutput()
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

func (s *ShellService) Shutdown(ctx context.Context) {
	select {
	case s.done <- struct{}{}:
	case <-ctx.Done():
	}
}
