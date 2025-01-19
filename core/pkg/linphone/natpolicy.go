package linphone

// #include <linphone/linphonecore.h>
import "C"
import "unsafe"

type NatPolicy C.LinphoneNatPolicy

func (n *NatPolicy) cPtr() *C.LinphoneNatPolicy {
	return (*C.LinphoneNatPolicy)(n)
}

func (n *NatPolicy) Ptr() *NatPolicy {
	return n
}

func (n *NatPolicy) SetStunServer(stunServer string) {
	s := C.CString(stunServer)
	defer C.free(unsafe.Pointer(s))

	C.linphone_nat_policy_set_stun_server(n.cPtr(), s)
}

func (n *NatPolicy) EnableStun(enable bool) {
	C.linphone_nat_policy_enable_stun(n.cPtr(), boolToCBool(enable))
}

func (n *NatPolicy) EnableIce(enable bool) {
	C.linphone_nat_policy_enable_stun(n.cPtr(), boolToCBool(enable))
}
