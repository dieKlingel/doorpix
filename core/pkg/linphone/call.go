package linphone

// #cgo darwin CFLAGS: -I/Applications/Linphone.app/Contents/Resources/include
// #cgo darwin LDFLAGS: -L/Applications/Linphone.app/Contents/Frameworks
// #cgo linux CFLAGS: -I/usr/include/linphone
// #cgo linux LDFLAGS: -llinphone
//
// #include <linphone/linphonecore.h>
import "C"

type Call C.LinphoneCall

func (call *Call) cPtr() *C.LinphoneCall {
	return (*C.LinphoneCall)(call)
}

func (call *Call) Accept() int {
	v := C.linphone_call_accept(call.cPtr())
	return int(v)
}

func (call *Call) SetSpeakerMuted(muted bool) {
	C.linphone_call_set_speaker_muted(call.cPtr(), boolToCBool(muted))
}
