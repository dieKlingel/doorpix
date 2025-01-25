package appvideo

// #include "driver.h"
import "C"
import (
	"log/slog"
	"unsafe"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
	"github.com/go-gst/go-gst/gst"
)

func Initialize() {
	if !pjsua2.EndpointInstance().LibIsThreadRegistered() {
		pjsua2.EndpointInstance().LibRegisterThread("appvideo_driver")
	}

	C.init_app_video()

}

func SetCameraFactory(factory func() *camera.Camera) {
	createNewCamera = factory
}

var streams = make(map[uintptr]Stream)
var createNewCamera func() *camera.Camera

func createNewCameraSafe() *camera.Camera {
	if createNewCamera == nil {
		slog.Warn("no camera is set for the appvideo driver")
		c, err := camera.NewFromString(
			"videotestsrc",
			camera.NewElement(
				"capsfilter",
				"caps", gst.NewCapsFromString("video/x-raw,format=I420,width=1920,height=1080,framerate=30/1"),
			),
		)
		if err != nil {
			panic(err)
		}

		return c
	}

	return createNewCamera()
}

//export go_video_stream_start
func go_video_stream_start(ptr unsafe.Pointer) {
	cam := createNewCameraSafe()
	err := cam.Start()
	if err != nil {
		slog.Error("failed to start stream", "error", err)
	}

	streams[uintptr(ptr)] = cam
}

//export go_video_stream_stop
func go_video_stream_stop(ptr unsafe.Pointer) {
	stream, ok := streams[uintptr(ptr)]
	if !ok {
		slog.Error("the stream should be closed but it does not exist", "id", uintptr(ptr))
		return
	}

	err := stream.Stop()
	if err != nil {
		slog.Error("failed to stop stream", "error", err)
	}
	delete(streams, uintptr(ptr))
}

//export go_video_stream_get_frame
func go_video_stream_get_frame(streamPtr *C.pjmedia_vid_dev_stream, framePtr *C.pjmedia_frame) {
	ptr := unsafe.Pointer(streamPtr)
	stream, ok := streams[uintptr(ptr)]
	if !ok {
		return
	}

	frame, ok := <-stream.Frame()
	if !ok {
		slog.Debug("frame channel closed")
	}
	expectedSize := int(framePtr.size)
	recivedSize := len(frame)

	usedSize := min(expectedSize, recivedSize)
	if expectedSize != recivedSize {
		slog.Warn("frame size mismatch", "expected", expectedSize, "recived", recivedSize)
	}

	bufferPtr := (*[3110400]C.uchar)(framePtr.buf)
	for i := 0; i < usedSize; i++ {
		bufferPtr[i] = (C.uchar)(frame[i])
	}
}
