package linphone

// #cgo darwin CFLAGS: -I/Applications/Linphone.app/Contents/Resources/include
// #cgo darwin LDFLAGS: -L/Applications/Linphone.app/Contents/Frameworks
// #cgo linux CFLAGS: -I/usr/include/linphone -I/usr/include/ortp
// #cgo linux LDFLAGS: -llinphone
//
// #include <linphone/linphonecore.h>
// // #include <ortp/payloadtype.h>
import "C"

type PayloadType C.LinphonePayloadType

func (p *PayloadType) cPtr() *C.LinphonePayloadType {
	return (*C.LinphonePayloadType)(p)
}

func (t *PayloadType) GetDescription() string {
	d := C.linphone_payload_type_get_description(t.cPtr())
	return C.GoString(d)
}
