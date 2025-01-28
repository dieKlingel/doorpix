package core

import (
	"context"
	"sync"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type Service interface{}

type InitService interface {
	Init() error
}

type RunnerService interface {
	Run(act doorpix.Action, event *doorpix.ActionHook) bool
}

type ExecService interface {
	Exec(ctx context.Context, wg *sync.WaitGroup) error
}
