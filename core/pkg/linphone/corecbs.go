package linphone

// #include <linphone/linphonecore.h>
//
// extern void global_state_changed(LinphoneCore *lc, LinphoneGlobalState gstate, char *message);
// extern void registration_state_changed(LinphoneCore *lc, LinphoneProxyConfig *cfg, LinphoneRegistrationState state, char *message);
// extern void call_state_changed(LinphoneCore *lc, LinphoneCall *call, LinphoneCallState gstate, char *message);
// extern void dtmf_received(LinphoneCore *lc, LinphoneCall *call, int dtmf);
import "C"
import (
	"log/slog"
	"unsafe"
)

type LinphoneCoreCbs C.LinphoneCoreCbs

type OnGlobalStateChangedCallback func(core *Core, state LinphoneGlobalState, message string)
type OnRegistrationStateChangedCallback func(core *Core, cfg *ProxyConfig, state LinphoneRegistrationState, message string)
type OnCallStateChangedCallback func(core *Core, call *Call, state CallState, message string)
type OnDtmfReceivedCallback func(core *Core, call *Call, dtmf int)

//export global_state_changed
func global_state_changed(c *C.LinphoneCore, state C.LinphoneGlobalState, message *C.char) {
	core := (*Core)(c)
	userdata := core.GetUserData()

	if userdata.onGlobalStateChanged != nil {
		userdata.onGlobalStateChanged(
			core,
			LinphoneGlobalState(state),
			C.GoString(message),
		)
	}
}

//export registration_state_changed
func registration_state_changed(c *C.LinphoneCore, proxycfg *C.LinphoneProxyConfig, state C.LinphoneRegistrationState, message *C.char) {
	slog.Debug("registration_state_changed", "core", c, "proxycfg", proxycfg, "state", state, "message", message)

	core := (*Core)(c)
	userdata := core.GetUserData()

	if userdata.onRegistrationStateChanged != nil {
		userdata.onRegistrationStateChanged(
			core,
			(*ProxyConfig)(proxycfg),
			LinphoneRegistrationState(state),
			C.GoString(message),
		)
	}
}

//export call_state_changed
func call_state_changed(c *C.LinphoneCore, call *C.LinphoneCall, state C.LinphoneCallState, message *C.char) {
	slog.Debug("call_state_changed", "core", c, "call", call, "state", state, "message", message)

	core := (*Core)(c)
	userdata := core.GetUserData()

	if userdata.onCallStateChanged != nil {
		userdata.onCallStateChanged(
			core,
			(*Call)(call),
			CallState(state),
			C.GoString(message),
		)
	}
}

//export dtmf_received
func dtmf_received(c *Core, e *C.LinphoneCall, dtmf C.int) {
	core := (*Core)(c)
	userdata := core.GetUserData()

	if userdata.onDtmfReceived != nil {
		userdata.onDtmfReceived(
			core,
			(*Call)(e),
			int(dtmf),
		)
	}
}

// Go Callbacks

func initCoreCbs(cbs *C.LinphoneCoreCbs) {
	C.linphone_core_cbs_set_global_state_changed(cbs, (*[0]byte)(unsafe.Pointer(C.global_state_changed)))
	C.linphone_core_cbs_set_registration_state_changed(cbs, (*[0]byte)(unsafe.Pointer(C.registration_state_changed)))
	C.linphone_core_cbs_set_call_state_changed(cbs, (*[0]byte)(unsafe.Pointer(C.call_state_changed)))
	C.linphone_core_cbs_set_dtmf_received(cbs, (*[0]byte)(unsafe.Pointer(C.dtmf_received)))
}

func (core *Core) OnGlobalStateChanged(f OnGlobalStateChangedCallback) {
	core.GetUserData().onGlobalStateChanged = f
}

func (core *Core) OnRegistrationStateChanged(f OnRegistrationStateChangedCallback) {
	core.GetUserData().onRegistrationStateChanged = f
}

func (core *Core) OnCallStateChanged(f OnCallStateChangedCallback) {
	core.GetUserData().onCallStateChanged = f
}
