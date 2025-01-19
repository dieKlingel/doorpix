package mediastreamer2

// #include <mediastreamer2/allfilters.h>
// #include <mediastreamer2/mscommon.h>
// #include <mediastreamer2/msjpegwriter.h>
// #include <mediastreamer2/msticker.h>
// #include <mediastreamer2/msvideo.h>
//#include <mediastreamer2/mswebcam.h>
import "C"

type Ticker = C.MSTicker

func NewTicker() *Ticker {
	return C.ms_ticker_new()
}

func (t *Ticker) Attach(f *Filter) int {
	r := C.ms_ticker_attach(t, f.cPtr())
	return int(r)
}

func (t *Ticker) Detach(f *Filter) int {
	r := C.ms_ticker_detach(t, f.cPtr())
	return int(r)
}

func (t *Ticker) Destroy() {
	C.ms_ticker_destroy(t)
}
