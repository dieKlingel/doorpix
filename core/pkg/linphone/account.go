package linphone

// #include <linphone/linphonecore.h>
import "C"

type Account C.LinphoneAccount

func (a *Account) cPtr() *C.LinphoneAccount {
	return (*C.LinphoneAccount)(a)
}
