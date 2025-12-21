package main

import (
	"context"
	"log/slog"
)

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}

	return value
}

type Server interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

func serve(srv Server) {
	go func() {
		err := srv.Serve()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
}

func shutdown(srv Server, ctx context.Context) {
	err := srv.Shutdown(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
}
