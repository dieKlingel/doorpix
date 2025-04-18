package actions_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCall struct {
	mock.Mock
}

func (m *MockCall) Hangup() error {
	args := m.Called()
	return args.Error(0)
}

func TestHangupActionExecute(t *testing.T) {
	t.Run("hangup call", func(t *testing.T) {
		action := &actions.HangupAction{}

		call := &MockCall{}
		call.On("Hangup").Return(nil)

		err := action.Execute(call, nil)
		assert.NoError(t, err)
		call.AssertExpectations(t)
	})

	t.Run("hangup call with error", func(t *testing.T) {
		action := &actions.HangupAction{}

		call := &MockCall{}
		call.On("Hangup").Return(assert.AnError)

		err := action.Execute(call, nil)
		assert.Error(t, err)
		call.AssertExpectations(t)
	})
}
