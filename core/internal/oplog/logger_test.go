package oplog_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/dieklingel/doorpix/core/internal/oplog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerPush(t *testing.T) {
	logger := oplog.NewLogger()

	t.Run("should not return error without verifier and listener", func(t *testing.T) {
		err := logger.Push("test message")
		assert.NoError(t, err)
	})

	t.Run("should return error with verifier", func(t *testing.T) {
		logger.AddVerifyFunc(func(msg any) error {
			return errors.New("verification failed")
		})

		err := logger.Push("test message")
		assert.Error(t, err)
	})

	t.Run("should not return error with non matching verifier", func(t *testing.T) {
		logger.AddVerifyFunc(func(msg any) error {
			_, ok := msg.(fmt.Stringer)
			if !ok {
				return errors.New("verification failed")
			}

			return nil
		})

		err := logger.Push("non matching message")
		assert.Error(t, err)
	})
}
