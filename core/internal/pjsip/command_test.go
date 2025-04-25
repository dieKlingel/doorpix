package pjsip_test

import (
	"testing"
	"time"

	"github.com/dieklingel/doorpix/core/internal/pjsip"
	"github.com/stretchr/testify/assert"
)

func TestSendInstantMessageCommand(t *testing.T) {
	channel := make(chan pjsip.Command)
	go func() {
		cmd := <-channel
		cmd.Error() <- nil
	}()

	cmd := pjsip.SendInstantMessageCommand{}
	channel <- &cmd
	select {
	case err := <-cmd.Error():
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("expected to receive an error, but timed out")
	}
}
