package mediastreamer2

// #include <mediastreamer2/allfilters.h>
// #include <mediastreamer2/mscommon.h>
// #include <mediastreamer2/msjpegwriter.h>
// #include <mediastreamer2/msticker.h>
// #include <mediastreamer2/msvideo.h>
//#include <mediastreamer2/mswebcam.h>
//
// extern void notify_callback(void *userdata, MSFilter *f, unsigned int id, void *arg);
import "C"
import "unsafe"

type Filter C.MSFilter

func (f *Filter) cPtr() *C.MSFilter {
	return (*C.MSFilter)(f)
}

func (f *Filter) ImplementsInterface(iface FilterInterfaceId) bool {
	return C.ms_filter_implements_interface(f.cPtr(), iface) != 0
}

func (f *Filter) CallMethod(m FilterMethod, p *any) int {
	param := unsafe.Pointer(p)
	defer C.free(param)

	r := C.ms_filter_call_method(f.cPtr(), C.uint(m), param)
	return int(r)
}

func (f *Filter) GetId() FilterId {
	return C.ms_filter_get_id(f.cPtr())
}

func (f *Filter) Link(p0 int, dst *Filter, p1 int) int {
	r := C.ms_filter_link(f.cPtr(), C.int(p0), dst.cPtr(), C.int(p1))
	return int(r)
}

func (f *Filter) Unlink(p0 int, dst *Filter, p1 int) int {
	r := C.ms_filter_unlink(f.cPtr(), C.int(p0), dst.cPtr(), C.int(p1))
	return int(r)
}

func (f *Filter) Destroy() {
	C.ms_filter_destroy(f.cPtr())
}

type callback struct {
	fun func(filter *Filter)
}

//export notify_callback
func notify_callback(userdata unsafe.Pointer, filter *Filter, id C.uint, arg unsafe.Pointer) {
	c := (*callback)(userdata)
	c.fun(filter)
}

func (f *Filter) AddNotifyCallback(fun func(filter *Filter)) {
	data := unsafe.Pointer(&callback{fun: fun})

	cb := unsafe.Pointer(C.notify_callback)
	C.ms_filter_add_notify_callback(f.cPtr(), (*[0]byte)(cb), data, 1)
}
