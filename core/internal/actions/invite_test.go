package actions_test

import (
	"testing"
	"text/template"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCaller struct {
	mock.Mock
}

func (m *MockCaller) Invite(uris []string) error {
	args := m.Called(uris)
	return args.Error(0)
}

func TestIviteActionExecute(t *testing.T) {
	t.Run("execute non existing command", func(t *testing.T) {
		action := &actions.InviteAction{
			UriTemplates: []template.Template{
				*template.Must(template.New("").Parse("sip:user@example.com")),
			},
		}
		caller := &MockCaller{}
		caller.On("Invite", []string{"sip:user@example.com"}).Return(nil)

		err := action.Execute(caller, map[string]any{})
		assert.NoError(t, err)
		caller.AssertExpectations(t)
	})
}
