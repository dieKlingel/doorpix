package actions_test

import (
	"testing"
	"text/template"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(topic string, message string) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

func TestPublishActionExecute(t *testing.T) {
	t.Run("publish message", func(t *testing.T) {
		action := &actions.PublishAction{
			Topic:   *template.Must(template.New("").Parse("test/topic")),
			Message: *template.Must(template.New("").Parse("Hello, World!")),
		}

		publisher := &MockPublisher{}
		publisher.On("Publish", "test/topic", "Hello, World!").Return(nil)

		err := action.Execute(publisher, nil)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		publisher.AssertExpectations(t)
	})

	t.Run("publish message with error", func(t *testing.T) {
		action := &actions.PublishAction{
			Topic:   *template.Must(template.New("").Parse("test/topic")),
			Message: *template.Must(template.New("").Parse("Hello, World!")),
		}

		publisher := &MockPublisher{}
		publisher.On("Publish", "test/topic", "Hello, World!").Return(assert.AnError)

		err := action.Execute(publisher, nil)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		publisher.AssertExpectations(t)
	})
}
