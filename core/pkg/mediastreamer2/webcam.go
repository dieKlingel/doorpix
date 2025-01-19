package mediastreamer2

// #include <mediastreamer2/allfilters.h>
// #include <mediastreamer2/mscommon.h>
// #include <mediastreamer2/msjpegwriter.h>
// #include <mediastreamer2/msticker.h>
// #include <mediastreamer2/msvideo.h>
//#include <mediastreamer2/mswebcam.h>
import "C"

type Webcam C.MSWebCam

func (w *Webcam) cPtr() *C.MSWebCam {
	return (*C.MSWebCam)(w)
}

func (w *Webcam) CreateReader() *Filter {
	return (*Filter)(C.ms_web_cam_create_reader(w.cPtr()))
}
