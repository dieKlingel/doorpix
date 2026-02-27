package camera

// #cgo pkg-config: --static libpjproject
// #include "driver.h"
import "C"
import (
	"log/slog"
	"unsafe"

	"github.com/dieklingel/doorpix/internal/media/camera"
)

const (
	PJ_SUCCESS = 0
	PJ_FAILURE = 1

	DRIVER_NAME = "github.com/dieklingel/doorpix"
)

type Webcam interface {
	Start() (camera.Session, error)
}

type CameraDriver struct {
	webcam   Webcam
	sessions map[uintptr]camera.Session
}

var driver *CameraDriver = nil

func Register(name string, webcam Webcam) error {
	if driver != nil {
		return ErrAlreadyInitialized
	}

	driver = &CameraDriver{
		webcam:   webcam,
		sessions: make(map[uintptr]camera.Session),
	}

	status := C.Register(C.factory_options{
		width:              720,
		height:             480,
		framerate:          15,
		name:               C.CString(name),
		name_length:        C.int(len(name)),
		driver_name:        C.CString(DRIVER_NAME),
		driver_name_length: C.int(len(DRIVER_NAME)),
	})

	if status == PJ_SUCCESS {
		return nil
	}

	return ErrRegistrationFailed
}

//export go_stream_start
func go_stream_start(ptr unsafe.Pointer) C.int {
	slog.Debug("start a new stream", "ptr", ptr, "driver", DRIVER_NAME)

	session, err := driver.webcam.Start()

	if err != nil {
		slog.Debug("driver: failed to start stream", "error", err.Error())
		return PJ_FAILURE
	}

	driver.sessions[uintptr(ptr)] = session
	return PJ_SUCCESS
}

//export go_stream_stop
func go_stream_stop(ptr unsafe.Pointer) C.int {
	slog.Debug("stop a camera stream", "ptr", ptr, "driver", DRIVER_NAME)

	session, exists := driver.sessions[uintptr(ptr)]
	if !exists {
		return PJ_SUCCESS
	}

	err := session.Stop()
	if err != nil {
		slog.Debug("driver: failed to stop stream", "error", err.Error())
		return PJ_FAILURE
	}

	delete(driver.sessions, uintptr(ptr))
	return PJ_SUCCESS
}

//export go_stream_get_frame
func go_stream_get_frame(streamPtr *C.pjmedia_vid_dev_stream, framePtr *C.pjmedia_frame) C.int {
	ptr := unsafe.Pointer(streamPtr)
	session, exists := driver.sessions[uintptr(ptr)]
	if !exists {
		return PJ_FAILURE
	}

	frame, ok := <-session.Frame()
	if !ok {
		slog.Debug("driver: failed to get frame from session, as the channel was closed")
		return PJ_FAILURE
	}

	expectedSize := int(framePtr.size)
	recivedSize := len(frame)

	usedSize := min(expectedSize, recivedSize)
	if expectedSize != recivedSize {
		slog.Warn("driver: video frame size mismatch", "expected", expectedSize, "recived", recivedSize)
	}

	// The buffer size has to be of a size, that the whole image fits in it.
	// The image is encoded in I420/YUV where 1,5 bytes per pixel are neccessary.
	// If the image has a size of 1920x1080 the buffer size has to be
	//   1920 * 1080 * 1,5 = 3.110.400 bytes
	bufferPtr := (*[3_110_400]C.uchar)(framePtr.buf)
	for i := range usedSize {
		bufferPtr[i] = (C.uchar)(frame[i])
	}
	return PJ_SUCCESS
}
