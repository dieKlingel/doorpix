package core

import (
	"bytes"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/exec"
)

type SystemService struct {
	Bus    *EventQueue
	Config doorpix.Config
}

func (service *SystemService) Name() string {
	return "system-service"
}

func (service *SystemService) Run(action doorpix.Action, hook *doorpix.ActionHook) bool {
	switch action := action.(type) {
	case doorpix.ConditionAction:
		var buf bytes.Buffer
		if err := action.Condition.Execute(&buf, hook); err != nil {
			slog.Error(err.Error())
			break
		}
		condition, ok := strconv.ParseBool(buf.String())
		if ok != nil {
			slog.Error(ok.Error())
		}
		if condition {
			hook.AdditionalActions = append(hook.AdditionalActions, action.Then...)
		} else {
			hook.AdditionalActions = append(hook.AdditionalActions, action.Else...)
		}

	case doorpix.LogAction:
		msg := bytes.Buffer{}
		if err := action.Message.Execute(&msg, hook.Data); err != nil {
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
			if err := expr.Execute(&buf, hook.Data); err != nil {
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

	default:
		return false
	}

	return true
}

func (service *SystemService) Init() error {
	slog.Debug("init system service")

	// TODO: init system service

	slog.Debug("successfully initialized system service")
	return nil
}
