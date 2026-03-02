package server

import "context"

type Worker interface {
	Run() error
	Stop(ctx context.Context) error
}
