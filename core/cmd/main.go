package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/env"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: env.LogLevel(),
	}))
	slog.SetDefault(logger)

	emitter := eventemitter.New()
	pipelineBuilder := core.NewPipelineBuilder()
	cfg := must(buildConfig())

	httpServer := core.NewHTTPServer(core.HTTPServerProps{
		Port:                    cfg.HTTP.Port,
		VideoStreamCameraDevice: cfg.Camera.Device,
	})
	if cfg.HTTP.Enabled {
		httpServer.Start()
	}

	grpcServer := core.NewGRPCServer(emitter, core.GRPCServerProps{
		Address: cfg.RPC.Host,
		Port:    cfg.RPC.Port,
	})
	if cfg.RPC.Enabled {
		grpcServer.Start()
	}

	gpioController := core.NewGPIOController(emitter, core.GPIOControllerProps{
		Pins: cfg.GPIO.Pins,
		Chip: cfg.GPIO.Chip,
	})
	if cfg.GPIO.Enabled {
		gpioController.Start()
	}

	mqttClient := core.NewMQTTClient(emitter, core.MQTTClientProps{
		ClientId:      cfg.MQTT.ClientId,
		Host:          cfg.MQTT.Host,
		Port:          cfg.MQTT.Port,
		Protocol:      cfg.MQTT.Protocol,
		Username:      cfg.MQTT.Username,
		Password:      cfg.MQTT.Password,
		Subscriptions: cfg.MQTT.Subscriptions,
	})
	if cfg.MQTT.Enabled {
		mqttClient.Start()
		pipelineBuilder.WithMqttMessagePublisher(mqttClient)
	}

	softphone := core.NewSIPClient(emitter, core.SIPClientProps{
		Username:    cfg.SIPPhone.Username,
		Password:    cfg.SIPPhone.Password,
		Realm:       cfg.SIPPhone.Realm,
		Server:      cfg.SIPPhone.Server,
		StunServers: cfg.SIPPhone.StunServers,
		VideoDevice: cfg.Camera.Device,
		Whitelist:   cfg.SIPPhone.Whitelist,
	})
	if cfg.SIPPhone.Enabled {
		if err := softphone.Start(); err != nil {
			slog.Error("failed to start SIP client", "error", err)
			os.Exit(1)
		}
		pipelineBuilder.WithSoftphone(softphone)
	}

	app := core.App{
		Config:              cfg,
		EventEmitter:        emitter,
		BasePipelineBuilder: pipelineBuilder,
	}
	if err := app.Start(); err != nil {
		panic(err)
	}

	emitter.Emit("events/startup", map[string]any{})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	slog.Info("Received interrupt signal, shutting down...")

	emitter.Emit("events/shutdown", map[string]any{})

	mqttClient.Stop()
	gpioController.Stop()
	grpcServer.Stop()
	httpServer.Stop()
}

func must[T any](val T, err error) T {
	if err != nil {
		slog.Error("failed to initialize", "error", err)
		os.Exit(1)
	}

	return val
}

func buildConfig() (doorpix.Config, error) {
	config := doorpix.NewConfig()
	config.AddConfigPath(
		"/etc/doorpix/doorpix.yaml",
		"/etc/doorpix/config.yaml",
		"doorpix.yaml",
		"config.yaml",
	)
	if err := config.Read(); err != nil {
		return doorpix.Config{}, err
	}
	if err := config.Error(); err != nil {
		return doorpix.Config{}, err
	}

	return config, nil
}
