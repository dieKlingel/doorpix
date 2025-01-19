package exec_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/exec"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	t.Run("run non existing command", func(t *testing.T) {
		_, err := exec.Run([]string{"command_which_does_not_exist"})
		assert.Error(t, err)
	})

	t.Run("return exit code 1", func(t *testing.T) {
		_, err := exec.Run([]string{"exit 1"})
		assert.Error(t, err)
	})

	t.Run("return exit code 0", func(t *testing.T) {
		_, err := exec.Run([]string{"exit 0"})
		assert.NoError(t, err)
	})
}
