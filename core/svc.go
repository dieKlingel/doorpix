package core

import (
	"context"
	"sync"
)

type Service interface{}

type InitService interface {
	Init() error
}

type ExecService interface {
	Exec(ctx context.Context, wg *sync.WaitGroup) error
}
