package linphone

// #include <linphone/linphonecore.h>
import "C"

type Transports C.LinphoneTransports

func (t *Transports) cPtr() *C.LinphoneTransports {
	return (*C.LinphoneTransports)(t)
}

func (t *Transports) SetUdpPort(port int) {
	p := C.int(port)
	C.linphone_transports_set_udp_port(t.cPtr(), p)
}

func (t *Transports) SetDtlsPort(port int) {
	p := C.int(port)
	C.linphone_transports_set_dtls_port(t.cPtr(), p)
}

func (t *Transports) SetTcpPort(port int) {
	p := C.int(port)
	C.linphone_transports_set_tcp_port(t.cPtr(), p)
}

func (t *Transports) SetTlsPort(port int) {
	p := C.int(port)
	C.linphone_transports_set_tls_port(t.cPtr(), p)
}
