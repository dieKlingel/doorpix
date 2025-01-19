package linphone

// #include <linphone/linphonecore.h>
import "C"
import "unsafe"

type Config C.LinphoneConfig

func (c *Config) cPtr() *C.LinphoneConfig {
	return (*C.LinphoneConfig)(c)
}

func (c *Config) SetBool(section string, key string, value bool) {
	s := C.CString(section)
	defer C.free(unsafe.Pointer(s))

	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))

	C.linphone_config_set_bool(c.cPtr(), s, k, boolToCBool(value))
}
