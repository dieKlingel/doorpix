package linphone

// #include <linphone/linphonecore.h>
import "C"

type LinphoneGlobalState C.LinphoneGlobalState

const (
	LinphoneGlobalStateOff         LinphoneGlobalState = C.LinphoneGlobalOff
	LinphoneGlobalStateStartup     LinphoneGlobalState = C.LinphoneGlobalStartup
	LinphoneGlobalStateOn          LinphoneGlobalState = C.LinphoneGlobalOn
	LinphoneGlobalStateShutdown    LinphoneGlobalState = C.LinphoneGlobalShutdown
	LinphoneGlobalStateConfiguring LinphoneGlobalState = C.LinphoneGlobalConfiguring
	LinphoneGlobalStateReady       LinphoneGlobalState = C.LinphoneGlobalReady
)
