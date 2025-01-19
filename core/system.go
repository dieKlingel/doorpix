package core

import (
	"bytes"
	"fmt"
	"log/slog"
	"time"

	"github.com/dieklingel/doorpix/core/internal/config"
	"github.com/dieklingel/doorpix/core/internal/exec"
)

type SystemHandler struct{}

func (s *SystemHandler) HandleEvent(action config.Action, event *Event) {
	switch action := action.(type) {
	case config.LogAction:
		msg := bytes.Buffer{}
		if err := action.Message.Execute(&msg, event.Data); err != nil {
			slog.Error(err.Error())
			break
		}

		slog.Info(msg.String())
	case config.SleepAction:
		time.Sleep(time.Duration(action.Duration) * time.Second)
	case config.EvalAction:
		out, err := exec.Run(action.Expressions)
		if err != nil {
			slog.Error(err.Error())
			break
		}
		fmt.Print(out)
	}
}

func (s *SystemHandler) Setup(emitter *EventEmitter) {
	emitter.Handler(s)
}

func (s *SystemHandler) Cleanup() {
}

func (s *SystemHandler) Exec() {
}
