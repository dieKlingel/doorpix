package mediastreamer2

// #include <mediastreamer2/allfilters.h>
// #include <mediastreamer2/mscommon.h>
// #include <mediastreamer2/msjpegwriter.h>
// #include <mediastreamer2/msticker.h>
// #include <mediastreamer2/msvideo.h>
// #include <mediastreamer2/mswebcam.h>
import "C"

type WebCamDesc C.MSWebCamDesc

func (desc *WebCamDesc) cPtr() *C.MSWebCamDesc {
	return (*C.MSWebCamDesc)(desc)
}
