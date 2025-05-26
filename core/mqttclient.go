package core

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/service"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type MQTTClientProps struct {
	ClientId      string
	Host          string
	Port          int
	Protocol      string
	Username      string
	Password      string
	Subscriptions []string
}

// TODO: implement last will message
type MQTTClient struct {
	props        MQTTClientProps
	eventemitter *eventemitter.EventEmitter

	ctx    service.Context
	client mqtt.Client
}

type MQTTMessageReceivedEvent struct {
	Topic   string
	Payload string
}

func NewMQTTClient(eventemitter *eventemitter.EventEmitter, props MQTTClientProps) *MQTTClient {
	return &MQTTClient{
		props:        props,
		eventemitter: eventemitter,

		ctx: service.NewContext(context.Background()),
	}
}

func (m *MQTTClient) Start() error {
	if len(m.props.Host) == 0 {
		return fmt.Errorf("mqtt host is required")
	}
	if len(m.props.Protocol) == 0 {
		return fmt.Errorf("mqtt protocol is required")
	}

	broker := fmt.Sprintf("%s://%s:%d", m.props.Protocol, m.props.Host, m.props.Port)
	clientId := m.props.ClientId
	if len(clientId) == 0 {
		clientId = fmt.Sprintf("doorpix-%s", uuid.NewString())
	}

	slog.Info("start the mqtt client", "broker", broker, "clientId", clientId)

	clientOptions := mqtt.NewClientOptions()
	clientOptions.
		AddBroker(broker).
		SetClientID(clientId).
		SetUsername(m.props.Username).
		SetPassword(m.props.Password).
		SetDefaultPublishHandler(m.onNewMessageReceived).
		SetResumeSubs(true).
		SetAutoReconnect(true)
	m.client = mqtt.NewClient(clientOptions)

	return m.exec()
}

func (m *MQTTClient) Stop() error {
	slog.Info("stop the mqtt client")
	return nil
}

func (m *MQTTClient) Publish(topic string, payload string) error {
	token := m.client.Publish(topic, 2, false, payload)
	sucess := token.WaitTimeout(1 * 1e9)
	if !sucess {
		return fmt.Errorf("failed to publish message to topic %s, error: %s", topic, token.Error().Error())
	}

	return nil
}

func (m *MQTTClient) exec() error {
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	for _, topic := range m.props.Subscriptions {
		slog.Debug("subscribing to mqtt topic", "topic", topic)

		token := m.client.Subscribe(topic, 2, nil)
		sucess := token.WaitTimeout(1 * time.Second)
		if !sucess {
			slog.Error("failed to subscribe to mqtt topic", "topic", topic, "error", token.Error().Error())
		}
	}

	m.ctx.Lock()
	go func() {
		defer m.ctx.Unlock()

		<-m.ctx.Done()
		slog.Debug("shutting down mqtt client")
		m.client.Disconnect(0)
		slog.Debug("successfully shut down mqtt client")
	}()

	return nil
}

func (m *MQTTClient) onNewMessageReceived(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	payload := string(message.Payload())

	slog.Debug("received new mqtt message", "topic", topic, "payload", payload)

	data := map[string]interface{}{
		"Topic":   topic,
		"Message": payload,
	}

	eventPath := fmt.Sprintf("events/mqtt-message/%s", topic)
	err := m.eventemitter.Emit(eventPath, data)
	if err != nil {
		slog.Error("failed to emit mqtt message received event", "error", err)
		return
	}
}
