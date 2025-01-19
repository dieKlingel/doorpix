package mediastreamer2

// #include <mediastreamer2/allfilters.h>
// #include <mediastreamer2/mscommon.h>
// #include <mediastreamer2/msjpegwriter.h>
// #include <mediastreamer2/msticker.h>
// #include <mediastreamer2/msvideo.h>
//#include <mediastreamer2/mswebcam.h>
import "C"

type Factory C.MSFactory

func (f *Factory) cPtr() *C.MSFactory {
	return (*C.MSFactory)(f)
}

func (f *Factory) CreateFilter(id FilterId) *Filter {
	return (*Filter)(C.ms_factory_create_filter(f.cPtr(), id))
}

func (f *Factory) GetWebCamManager() *WebCamManager {
	return (*WebCamManager)(C.ms_factory_get_web_cam_manager(f.cPtr()))
}
