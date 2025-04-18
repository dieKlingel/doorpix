package actions_test

import (
	"testing"
	"text/template"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/stretchr/testify/mock"
)

type MockMessanger struct {
	mock.Mock
}

func (m *MockMessanger) SendMessage(uris []string, message string) error {
	args := m.Called(uris, message)
	return args.Error(0)
}

func TestMessageExecute(t *testing.T) {
	t.Run("send message", func(t *testing.T) {
		action := &actions.MessageAction{
			UriTemplates: []template.Template{
				*template.Must(template.New("").Parse("sip:user@example.com")),
			},
			MessageTemplate: *template.Must(template.New("").Parse("Hello, World!")),
		}

		messanger := &MockMessanger{}
		messanger.On("SendMessage", []string{"sip:user@example.com"}, "Hello, World!").Return(nil)

		action.Execute(messanger, nil)
		messanger.AssertExpectations(t)
	})
}
