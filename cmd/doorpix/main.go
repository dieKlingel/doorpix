package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/dieklingel/doorpix/internal/app"
	"github.com/dieklingel/doorpix/internal/config"
)

func main() {
	cfg, err := config.NewBuilder().
		AddConfigFile("doorpix.yaml").
		AddConfigFile("~/doorpix.yaml").
		Build()

	if err != nil {
		slog.Error("failed to build config", "error", err)
		os.Exit(1)
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("starting doorpix...")

	cameraDriver := app.CreateCameraDriver(cfg)
	userAgent := app.CreateUserAgent(cfg, cameraDriver)
	httpServer := app.CreateHTTPServer(cfg, cameraDriver, userAgent)

	if httpServer != nil {
		serve(httpServer)
	}
	if userAgent != nil {
		serve(userAgent)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	slog.Info("stopping doorpix")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	if httpServer != nil {
		shutdown(httpServer, ctx)
	}
	if userAgent != nil {
		shutdown(userAgent, ctx)
	}
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
