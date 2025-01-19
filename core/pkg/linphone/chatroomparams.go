package linphone

// #include <linphone/linphonecore.h>
import "C"

type ChatRoomParams C.LinphoneChatRoomParams

func (c *ChatRoomParams) cPtr() *C.LinphoneChatRoomParams {
	return (*C.LinphoneChatRoomParams)(c)
}
