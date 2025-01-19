package linphone

// #include <linphone/linphonecore.h>
import "C"
import "unsafe"

type Factory C.LinphoneFactory

func (factory *Factory) Ptr() *Factory {
	return factory
}

func (factory *Factory) cPtr() *C.LinphoneFactory {
	return (*C.LinphoneFactory)(factory)
}

func GetFactory() *Factory {
	return (*Factory)(C.linphone_factory_get())
}

func CleanFactory() {
	C.linphone_factory_clean()
}

func (f *Factory) CreateCoreCbs() *C.LinphoneCoreCbs {
	cbs := C.linphone_factory_create_core_cbs(f.cPtr())

	initCoreCbs(cbs)

	return cbs
}

func (f *Factory) CreateLoggingServiceCbs() *LoggingServiceCbs {
	cbs := C.linphone_factory_create_logging_service_cbs(f.cPtr())
	return (*LoggingServiceCbs)(cbs)
}

func (f *Factory) CreateCore() *Core {
	core := (*Core)(C.linphone_factory_create_core_3(f.cPtr(), nil, nil, nil))
	core.SetUserData(&CoreUserData{})

	cbs := f.CreateCoreCbs()
	core.AddCallbacks(cbs)

	return core
}

func (f *Factory) CreateAddress(address string) *Address {
	addr := C.CString(address)
	defer C.free(unsafe.Pointer(addr))

	return (*Address)(C.linphone_factory_create_address(f.cPtr(), addr))
}

func (f *Factory) CreateAuthInfo(username string, userid string, password string, ha1 string, realm string, domain string) *AuthInfo {
	u := C.CString(username)
	defer C.free(unsafe.Pointer(u))

	i := C.CString(userid)
	defer C.free(unsafe.Pointer(i))

	p := C.CString(password)
	defer C.free(unsafe.Pointer(p))

	h := C.CString(ha1)
	defer C.free(unsafe.Pointer(h))

	r := C.CString(realm)
	defer C.free(unsafe.Pointer(r))

	d := C.CString(domain)
	defer C.free(unsafe.Pointer(d))

	return (*AuthInfo)(C.linphone_factory_create_auth_info(f.cPtr(), u, i, p, h, r, d))
}

func (f *Factory) CreateTransports() *Transports {
	return (*Transports)(C.linphone_factory_create_transports(f.cPtr()))
}

func (factory *Factory) SetDataResourcesDir(dir string) {
	d := C.CString(dir)
	defer C.free(unsafe.Pointer(d))

	C.linphone_factory_set_data_resources_dir(factory.cPtr(), d)
}

func (factory *Factory) SetImageResourcesDir(dir string) {
	d := C.CString(dir)
	defer C.free(unsafe.Pointer(d))

	C.linphone_factory_set_image_resources_dir(factory.cPtr(), d)
}

func (factory *Factory) SetTopResourcesDir(dir string) {
	d := C.CString(dir)
	defer C.free(unsafe.Pointer(d))

	C.linphone_factory_set_top_resources_dir(factory.cPtr(), d)
}
