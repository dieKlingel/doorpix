package linphone

// #include <linphone/linphonecore.h>
import "C"
import (
	"unsafe"

	"github.com/dieklingel/doorpix/core/pkg/mediastreamer2"
)

type Core C.LinphoneCore

func (c *Core) cPtr() *C.LinphoneCore {
	return (*C.LinphoneCore)(c)
}

func GetVersion() string {
	v := C.linphone_core_get_version()
	return C.GoString(v)
}

func (c *Core) GetMsFactory() *mediastreamer2.Factory {
	facptr := unsafe.Pointer(C.linphone_core_get_ms_factory(c.cPtr()))
	return (*mediastreamer2.Factory)(facptr)
}

func (c *Core) ReloadVideoDevices() {
	C.linphone_core_reload_video_devices(c.cPtr())
}

func (c *Core) SetVideoDevice(device string) {
	dev := C.CString(device)
	defer C.free(unsafe.Pointer(dev))

	C.linphone_core_set_video_device(c.cPtr(), dev)
}

func (c *Core) CreateChatRoom(params *ChatRoomParams, localAddr *Address, participants []*Address) *ChatRoom {
	if len(participants) == 0 {
		return nil
	}

	list := C.bctbx_list_new(unsafe.Pointer(participants[0].cPtr()))
	for _, p := range participants[1:] {
		list = C.bctbx_list_append(list, unsafe.Pointer(p.cPtr()))
	}

	room := C.linphone_core_create_chat_room_6(c.cPtr(), params.cPtr(), localAddr.cPtr(), list)
	return (*ChatRoom)(room)
}

func (c *Core) CreateDefaultChatRoomParams() *ChatRoomParams {
	p := C.linphone_core_create_default_chat_room_params(c.cPtr())
	return (*ChatRoomParams)(p)
}

func (c *Core) EnableIPv6(enabled bool) {
	C.linphone_core_enable_ipv6(c.cPtr(), boolToCBool(enabled))
}

func (c *Core) EnableSelfView(enabled bool) {
	C.linphone_core_enable_self_view(c.cPtr(), boolToCBool(enabled))
}

func (c *Core) EnableVideoCapture(enabled bool) {
	C.linphone_core_enable_video_capture(c.cPtr(), boolToCBool(enabled))
}

func (c *Core) EnableVideoDisplay(enabled bool) {
	C.linphone_core_enable_video_display(c.cPtr(), boolToCBool(enabled))
}

func (c *Core) CreateNatPolicy() *NatPolicy {
	return (*NatPolicy)(C.linphone_core_create_nat_policy(c.cPtr()))
}

func (c *Core) GetConfig() *Config {
	return (*Config)(C.linphone_core_get_config(c.cPtr()))
}

func (c *Core) AddCallbacks(l *C.LinphoneCoreCbs) {
	C.linphone_core_add_callbacks(c.cPtr(), l)
}

func (c *Core) SetTransports(t *Transports) {
	C.linphone_core_set_transports(c.cPtr(), t.cPtr())
}

func (c *Core) GetCurrentCall() *Call {
	return (*Call)(C.linphone_core_get_current_call(c.cPtr()))
}

func (c *Core) CreateCallParams(call *Call) *CallParams {
	return (*CallParams)(C.linphone_core_create_call_params(c.cPtr(), call.cPtr()))
}

func (c *Core) InviteWithParams(number string, params *CallParams) *Call {
	address := C.CString(number)
	defer C.free(unsafe.Pointer(address))

	return (*Call)(C.linphone_core_invite_with_params(c.cPtr(), address, params.cPtr()))
}

func (c *Core) Start() int {
	v := C.linphone_core_start(c.cPtr())
	return int(v)
}

func (c *Core) Stop() {
	C.linphone_core_stop(c.cPtr())
}

func (c *Core) RefreshRegisters() {
	C.linphone_core_refresh_registers(c.cPtr())
}

func (c *Core) Iterate() {
	C.linphone_core_iterate(c.cPtr())
}

func (c *Core) SetNatPolicy(n *NatPolicy) {
	C.linphone_core_set_nat_policy(c.cPtr(), n.Ptr().cPtr())
}

func (c *Core) SetUserData(data *CoreUserData) {
	p := unsafe.Pointer(data)
	C.linphone_core_set_user_data(c.cPtr(), p)
}

func (c *Core) GetUserData() *CoreUserData {
	u := C.linphone_core_get_user_data(c.cPtr())
	return (*CoreUserData)(u)
}

func (c *Core) AddAuthInfo(i *AuthInfo) {
	C.linphone_core_add_auth_info(c.cPtr(), i.cPtr())
}

func (c *Core) CreateAccountParams() *AccountParams {
	return (*AccountParams)(C.linphone_core_create_account_params(c.cPtr()))
}

func (c *Core) CreateAccount(p *AccountParams) *Account {
	return (*Account)(C.linphone_core_create_account(c.cPtr(), p.cPtr()))
}

func (c *Core) AddAccount(a *Account) {
	C.linphone_core_add_account(c.cPtr(), (*C.LinphoneAccount)(a.cPtr()))
}

func (c *Core) SetDefaultAccount(a *Account) {
	C.linphone_core_set_default_account(c.cPtr(), (*C.LinphoneAccount)(a.cPtr()))
}

func (c *Core) GetAudioPayloadTypes() []*PayloadType {
	types := C.linphone_core_get_audio_payload_types(c.cPtr())

	size := C.bctbx_list_size(types)
	items := make([]*PayloadType, int(size))

	item := types
	for i := 0; i < int(size); i++ {
		payload := (*PayloadType)(item.data)
		items[i] = payload
		item = item.next
	}

	return items
}

func (c *Core) GetVideoPayloadTypes() []*PayloadType {
	types := C.linphone_core_get_video_payload_types(c.cPtr())

	size := C.bctbx_list_size(types)
	items := make([]*PayloadType, int(size))

	item := types
	for i := 0; i < int(size); i++ {
		payload := (*PayloadType)(item.data)
		items[i] = payload
		item = item.next
	}

	return items
}

func (c *Core) GetGlobalState() LinphoneGlobalState {
	return (LinphoneGlobalState)(C.linphone_core_get_global_state(c.cPtr()))
}

func (c *Core) GetVideoDevicesList() []string {
	devices := C.linphone_core_get_video_devices_list(c.cPtr())

	size := C.bctbx_list_size(devices)

	//s := C.GoString(x)
	//println(s)

	items := make([]string, int(size))

	item := devices
	for i := 0; i < int(size); i++ {
		chars := (*C.char)(item.data)
		items[i] = C.GoString(chars)
		item = item.next
	}

	return items
}
