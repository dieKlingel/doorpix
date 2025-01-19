package mediastreamer2

// #include <mediastreamer2/allfilters.h>
// #include <mediastreamer2/mscommon.h>
// #include <mediastreamer2/msjpegwriter.h>
// #include <mediastreamer2/msticker.h>
// #include <mediastreamer2/msvideo.h>
//#include <mediastreamer2/mswebcam.h>
import "C"
import "unsafe"

type WebCamManager C.MSWebCamManager

func (m *WebCamManager) cPtr() *C.MSWebCamManager {
	return (*C.MSWebCamManager)(m)
}

func (m *WebCamManager) GetCam(device string) *Webcam {
	d := C.CString(device)
	defer C.free(unsafe.Pointer(d))

	return (*Webcam)(C.ms_web_cam_manager_get_cam(m.cPtr(), d))
}

func (m *WebCamManager) RegisterDesc(description *WebCamDesc) {
	C.ms_web_cam_manager_register_desc(m.cPtr(), description.cPtr())
}
