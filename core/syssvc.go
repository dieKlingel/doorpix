package core

import (
	"bytes"
	"fmt"
	"log/slog"
	"time"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/exec"
)

type SystemService struct {
	System doorpix.System
}

func (service *SystemService) HandleEvent(action doorpix.Action, event *doorpix.Event) {
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
		expressions := make([]string, len(action.Expressions))
		for i, expr := range action.Expressions {
			var buf bytes.Buffer
			if err := expr.Execute(&buf, event.Data); err != nil {
				slog.Error(err.Error())
				continue
			}
			expressions[i] = buf.String()
		}

		out, err := exec.Run(expressions)
		if err != nil {
			slog.Error(err.Error())
			break
		}
		fmt.Print(out)
	}
}

func (service *SystemService) Init() error {
	slog.Debug("init system service")

	service.System.Bus.Handler(service)

	slog.Debug("successfully initialized system service")
	return nil
}
