package shell_test

import (
	"testing"

	"github.com/dieklingel/doorpix/internal/device/shell"
	"github.com/stretchr/testify/assert"
)

func TestController_Exec(t *testing.T) {
	t.Run("should execute echo", func(t *testing.T) {
		controller := shell.NewController()
		out, err := controller.Exec("sh", "-c", "echo -n Hello World")

		assert.NoError(t, err)
		assert.Equal(t, "Hello World", string(out))
	})

	t.Run("should return error", func(t *testing.T) {
		controller := shell.NewController()
		_, err := controller.Exec("sh", "-c", "notABinary")

		assert.Error(t, err)
	})
}
