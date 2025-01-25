package appvideo

// #include "driver.h"
import "C"
import (
	"log/slog"
	"unsafe"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/go-gst/go-gst/gst"
)

func Initialize() {
	/*if !pjsua2.EndpointInstance().LibIsThreadRegistered() {
		pjsua2.EndpointInstance().LibRegisterThread("appvideoCameraDriverInitalization")
	}*/

	if !isInitalized {
		isInitalized = true
		C.init_app_video()
	}
}

func SetCameraDevice(dev string) {
	device = dev
	Initialize()
}

// The name of the camera device that is used by pjsip by this driver
func GetCameraDeviceName() string {
	return "DoorPiX Emulated Video Device"
}

var isInitalized bool = false
var streams = make(map[uintptr]Stream)
var device string = "videotestsrc"

func createNewCamera() *camera.Camera {
	c, err := camera.NewFromString(
		device,
		camera.MustNewElement("videoscale"),
		camera.MustNewElement(
			"capsfilter",
			"caps", gst.NewCapsFromString("video/x-raw,width=1920,height=1080"),
		),
		camera.MustNewElement("videoconvert"),
		camera.MustNewElement(
			"capsfilter",
			"caps", gst.NewCapsFromString("video/x-raw,format=I420,framerate=30/1"),
		),
	)

	if err != nil {
		panic(err)
	}

	return c
}

//export go_video_stream_start
func go_video_stream_start(ptr unsafe.Pointer) {
	cam := createNewCamera()
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
		slog.Warn("video frame size mismatch", "expected", expectedSize, "recived", recivedSize)
	}

	bufferPtr := (*[3110400]C.uchar)(framePtr.buf)
	for i := 0; i < usedSize; i++ {
		bufferPtr[i] = (C.uchar)(frame[i])
	}
}
