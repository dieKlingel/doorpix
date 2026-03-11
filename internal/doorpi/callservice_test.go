package doorpi_test

import (
	"context"
	"testing"

	"github.com/dieklingel/doorpix/internal/doorpi"
	"github.com/dieklingel/doorpix/internal/oplog"
	"github.com/dieklingel/doorpix/internal/transport/sip"
	"github.com/stretchr/testify/assert"
)

func TestCallService_OnCallEvent(t *testing.T) {
	t.Run("should call invite", func(t *testing.T) {
		ua := &MockUserAgent{}
		service := doorpi.NewSipService(ua)

		service.Listen()
		go service.Serve()

		ua.On("Invite", "sip:test@example.com").Return(&sip.CallInfo{}, nil)
		oplog.Dispatch("internal/doorpix/service/call/invite", "uri", "sip:test@example.com")

		err := service.Stop(context.Background())
		assert.NoError(t, err)

		ua.AssertExpectations(t)
	})
}

func TestCallService_OnMessageEvent(t *testing.T) {
	t.Run("should call message send", func(t *testing.T) {
		ua := &MockUserAgent{}
		service := doorpi.NewSipService(ua)

		service.Listen()
		go service.Serve()

		ua.On("SendMessage", "sip:test@example.com", "Hello World").Return(nil)
		oplog.Dispatch("internal/doorpix/service/call/message", "uri", "sip:test@example.com", "body", "Hello World")

		err := service.Stop(context.Background())
		assert.NoError(t, err)

		ua.AssertExpectations(t)
	})
}
