package core

import "context"

type MQTTClient struct {
}

func NewMQTTClient() *MQTTClient {
	return &MQTTClient{}
}

func (m *MQTTClient) Start(ctx context.Context) error {
	return nil
}

func (m *MQTTClient) Stop(ctx context.Context) error {
	return nil
}

func (m *MQTTClient) Publish(topic string, payload string) error {
	return nil
}
