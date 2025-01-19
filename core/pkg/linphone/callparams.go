package linphone

// #cgo darwin CFLAGS: -I/Applications/Linphone.app/Contents/Resources/include
// #cgo darwin LDFLAGS: -L/Applications/Linphone.app/Contents/Frameworks
// #cgo linux CFLAGS: -I/usr/include/linphone
// #cgo linux LDFLAGS: -llinphone
//
// #include <linphone/linphonecore.h>
import "C"

type CallParams C.LinphoneCallParams

func (c *CallParams) cPtr() *C.LinphoneCallParams {
	return (*C.LinphoneCallParams)(c)
}

func (c *CallParams) EnableEarlyMediaSending(enable bool) {
	C.linphone_call_params_enable_early_media_sending(c.cPtr(), boolToCBool(enable))
}

func (c *CallParams) SetMediaEncryption(mediaEncryption MediaEncryption) {
	C.linphone_call_params_set_media_encryption(c.cPtr(), mediaEncryption)
}
