package core

import (
	"context"
	"sync"
)

type Service interface {
	Name() string
}

type InitService interface {
	Init() error
}

type BackgroundService interface {
	StartBackgroundTask(ctx context.Context, wg *sync.WaitGroup) error
}
