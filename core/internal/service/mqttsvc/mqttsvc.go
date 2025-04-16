package mqttsvc

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTServiceProps struct {
	Host          string
	Port          int
	Protocol      string
	Username      string
	Password      string
	Subscriptions []string
}

type MQTTService struct {
	props MQTTServiceProps
	Emit  doorpix.Emit

	client mqtt.Client
}

func New(props MQTTServiceProps) *MQTTService {
	return &MQTTService{
		props: props,
	}
}

func (service *MQTTService) Name() string {
	return "mqtt-service"
}

func (service *MQTTService) Init() error {
	slog.Debug("init mqtt service")

	if len(service.props.Host) == 0 {
		return fmt.Errorf("mqtt host is required")
	}
	if len(service.props.Protocol) == 0 {
		return fmt.Errorf("mqtt protocol is required")
	}

	broker := fmt.Sprintf("%s://%s:%d", service.props.Protocol, service.props.Host, service.props.Port)
	slog.Debug("connecting to mqtt broker", "broker", broker)

	options := mqtt.NewClientOptions()
	options.AddBroker(broker)
	options.AutoReconnect = true
	options.SetUsername(service.props.Username)
	options.SetPassword(service.props.Password)
	options.SetDefaultPublishHandler(service.onNewMessageReceived)
	options.OnConnect = service.onConnectionEstablished
	options.OnConnectionLost = service.onConnectionLost
	service.client = mqtt.NewClient(options)

	slog.Debug("successfully initialized mqtt service")
	return nil
}

func (service *MQTTService) Publish(topic string, payload string) error {
	token := service.client.Publish(topic, 2, false, payload)

	sucess := token.WaitTimeout(1 * time.Second)
	if !sucess {
		return fmt.Errorf("failed to publish mqtt message on topic: %s", topic)
	}

	return nil
}

func (service *MQTTService) Run(act doorpix.Action, hook *doorpix.ActionHook) bool {
	action, ok := act.(doorpix.PublishAction)
	if !ok {
		return false
	}
	slog.Debug("run mqtt service")

	var topic bytes.Buffer
	err := action.Topic.Execute(&topic, hook.Data)
	if err != nil {
		slog.Error("failed to execute topic template", "error", err)
		return false
	}

	var payload bytes.Buffer
	err = action.Message.Execute(&payload, hook.Data)
	if err != nil {
		slog.Error("failed to execute message template", "error", err)
		return false
	}

	token := service.client.Publish(topic.String(), 2, false, payload.String())
	sucess := token.WaitTimeout(1 * time.Second)
	if !sucess {
		slog.Error("failed to publish mqtt message", "topic", topic.String())
		return false
	}

	return true
}

func (service *MQTTService) StartBackgroundTask(ctx context.Context, wg *sync.WaitGroup) error {
	slog.Debug("run mqtt service in background")

	token := service.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, topic := range service.props.Subscriptions {
		token := service.client.Subscribe(topic, 2, nil)
		sucess := token.WaitTimeout(1 * time.Second)
		if !sucess {
			slog.Error("failed to subscribe to mqtt topic", "topic", topic)
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		slog.Debug("shutting down mqtt service")
		service.client.Disconnect(0)
		slog.Debug("successfully shut down mqtt service")
	}()

	return nil
}

func (service *MQTTService) onNewMessageReceived(client mqtt.Client, msg mqtt.Message) {
	slog.Debug("received new mqtt message", "topic", msg.Topic(), "payload", string(msg.Payload()))

	service.Emit(
		doorpix.MqttMessageEvent,
		map[string]any{
			"topic":   msg.Topic(),
			"message": string(msg.Payload()),
		},
	)
}

func (service *MQTTService) onConnectionLost(client mqtt.Client, err error) {
	slog.Warn("mqtt connection lost", "error", err)
}

func (service *MQTTService) onConnectionEstablished(client mqtt.Client) {
	slog.Debug("mqtt connection established")
}
