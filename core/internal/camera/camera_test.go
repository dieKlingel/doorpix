package camera_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/go-gst/go-gst/gst"
	"github.com/stretchr/testify/assert"
)

func TestCamera_New(t *testing.T) {
	t.Run("create new videotestsrc", func(t *testing.T) {
		hardwareCamera, err := camera.NewHardwareCamera("videotestsrc")
		assert.NoError(t, err)

		cam, err := camera.NewFromHardwareCamera(hardwareCamera)
		assert.NoError(t, err)

		err = cam.Start()
		assert.NoError(t, err)

		_, ok := <-cam.Frame()
		assert.True(t, ok)

		err = cam.Stop()
		assert.NoError(t, err)

		_, ok = <-cam.Frame()
		assert.False(t, ok)
	})

	t.Run("create new jpeg encoded camera", func(t *testing.T) {
		webcam, err := camera.New("videotestsrc", camera.JPEG)
		assert.NoError(t, err)

		err = webcam.Start()
		assert.NoError(t, err)

		for i := 0; i < 100; i++ {
			_, ok := <-webcam.Frame()
			assert.True(t, ok)
		}

		err = webcam.Stop()
		assert.NoError(t, err)
	})
}

func TestCamera_ReadSingleFrame(t *testing.T) {
	cam, err := camera.NewFromString("videotestsrc",
		camera.MustNewElement("capsfilter", "caps", gst.NewCapsFromString("video/x-raw,format=I420,width=640,height=480")),
	)
	assert.NoError(t, err)

	cam.Start()
	frame, ok := <-cam.Frame()
	assert.True(t, ok)
	cam.Stop()

	assert.NotNil(t, frame)
	assert.Len(t, frame, 460800)
}
