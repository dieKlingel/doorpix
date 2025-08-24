package oplog_test

import (
	"fmt"
	"testing"

	"github.com/dieklingel/doorpix/core/internal/oplog"
	"github.com/dieklingel/doorpix/core/internal/workflow/event"
	"github.com/stretchr/testify/assert"
)

func TestListnerPush(t *testing.T) {

	t.Run("with primitive type", func(t *testing.T) {
		listener := &oplog.Listner[string]{
			Buffer: 10,
		}

		t.Run("valid message", func(t *testing.T) {
			success := listener.Push("test")
			assert.True(t, success)
		})

		t.Run("invalid message", func(t *testing.T) {
			success := listener.Push(123)
			assert.False(t, success)
		})
	})

	t.Run("with interface type", func(t *testing.T) {
		listener := &oplog.Listner[fmt.Stringer]{
			Buffer: 10,
		}

		t.Run("valid message", func(t *testing.T) {
			success := listener.Push(event.New("i am a stringer"))
			assert.True(t, success)
		})

		t.Run("invalid message", func(t *testing.T) {
			success := listener.Push(123)
			assert.False(t, success)
		})
	})
}

func TestListnerFilterFunc(t *testing.T) {
	listener := &oplog.Listner[string]{
		Buffer: 1,
	}

	t.Run("filters empty strings should success with non empty string", func(t *testing.T) {
		listener.FilterFunc(func(msg string) bool {
			return msg != ""
		})

		success := listener.Push("i am a stringer")
		assert.True(t, success)
	})

	t.Run("filters empty strings should fail with  empty string", func(t *testing.T) {
		listener.FilterFunc(func(msg string) bool {
			return msg != ""
		})

		success := listener.Push("")
		assert.False(t, success)
	})
}

func TestListnerListen(t *testing.T) {
	listener := &oplog.Listner[string]{
		Buffer: 1,
	}

	t.Run("valid message", func(t *testing.T) {
		value := "a message"
		listener.Push(value)
		val := <-listener.Listen()
		assert.Equal(t, value, val)
	})
}
