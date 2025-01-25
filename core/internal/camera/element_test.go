package camera_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/stretchr/testify/assert"
)

func TestNewElement(t *testing.T) {
	t.Run("create testvideosrc element", func(t *testing.T) {
		element := camera.MustNewElement("videotestsrc")
		assert.NotNil(t, element)
	})

	t.Run("create capsfilter element with properties", func(t *testing.T) {
		element := camera.MustNewElement("capsfilter", "caps", "video/x-raw,format=RGB")

		assert.NotNil(t, element)
	})

	t.Run("create capsfilter produces frame", func(t *testing.T) {

	})
}
