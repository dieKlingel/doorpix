package linphone

// #cgo darwin CFLAGS: -I/Applications/Linphone.app/Contents/Resources/include
// #cgo darwin LDFLAGS: -L/Applications/Linphone.app/Contents/Frameworks
// #cgo linux CFLAGS: -I/usr/include/linphone
// #cgo linux LDFLAGS: -llinphone
//
// #include <linphone/linphonecore.h>
import "C"

type AuthInfo C.LinphoneAuthInfo

func (authInfo *AuthInfo) cPtr() *C.LinphoneAuthInfo {
	return (*C.LinphoneAuthInfo)(authInfo)
}
