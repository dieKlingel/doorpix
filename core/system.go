package core

import (
	"bytes"
	"fmt"
	"log/slog"
	"time"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/exec"
)

type SystemHandler struct {
	System doorpix.System
}

func (s *SystemHandler) HandleEvent(action doorpix.Action, event *doorpix.Event) {
	switch action := action.(type) {
	case doorpix.LogAction:
		msg := bytes.Buffer{}
		if err := action.Message.Execute(&msg, event.Data); err != nil {
			slog.Error(err.Error())
			break
		}

		slog.Info(msg.String())
	case doorpix.SleepAction:
		time.Sleep(time.Duration(action.Duration) * time.Second)
	case doorpix.EvalAction:
		out, err := exec.Run(action.Expressions)
		if err != nil {
			slog.Error(err.Error())
			break
		}
		fmt.Print(out)
	}
}

func (s *SystemHandler) Setup() {
	s.System.Bus.Handler(s)
}

func (s *SystemHandler) Cleanup() {
}

func (s *SystemHandler) Exec() {
}
