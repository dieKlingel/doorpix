package camera_test

import (
	"testing"

	"github.com/dieklingel/doorpix/internal/media/camera"
	"github.com/stretchr/testify/assert"
)

func TestNewGstDriver(t *testing.T) {
	t.Run("should not return error for valid pipeline", func(t *testing.T) {
		driver, err := camera.NewGstDriver(`
			videotestsrc pattern=ball ! video/x-raw,width=640,height=480,framerate=30/1 ! tee name=tee
				tee. ! queue ! valve name=valve-test ! jpegenc ! appsink name=appsink-test
		`)

		assert.NotNil(t, driver)
		assert.NoError(t, err)
	})

	t.Run("should return error for invalid pipeline", func(t *testing.T) {
		driver, err := camera.NewGstDriver("invalidsrc")

		assert.Nil(t, driver)
		assert.Error(t, err)
	})
}
