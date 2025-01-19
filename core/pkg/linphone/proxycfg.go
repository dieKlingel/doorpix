package linphone

// #include <linphone/linphonecore.h>
import "C"

type ProxyConfig C.LinphoneProxyConfig

func (proxyConfig *ProxyConfig) Ptr() *ProxyConfig {
	return proxyConfig
}

func (proxyConfig *ProxyConfig) cPtr() *C.LinphoneProxyConfig {
	return (*C.LinphoneProxyConfig)(proxyConfig)
}
