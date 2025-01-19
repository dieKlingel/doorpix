package linphone

// #cgo linux CFLAGS: -I/usr/include/linphone
// #cgo linux LDFLAGS: -llinphone
//
// #include <linphone/linphonecore.h>
import "C"

type Address C.LinphoneAddress

func (address *Address) cPtr() *C.LinphoneAddress {
	return (*C.LinphoneAddress)(address)
}

func (address *Address) SetTransport(transport TransportType) {
	C.linphone_address_set_transport(address.cPtr(), *transport.cPtr())
}
