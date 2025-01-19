package linphone

// #include <linphone/linphonecore.h>
import "C"

type ChatRoom C.LinphoneChatRoom

func (c *ChatRoom) Ptr() *ChatRoom {
	return c
}

func (c *ChatRoom) cPtr() *C.LinphoneChatRoom {
	return (*C.LinphoneChatRoom)(c)
}

func (c *ChatRoom) CreateMessageFromUtf8(message string) *ChatMessage {
	msg := C.linphone_chat_room_create_message_from_utf8(c.cPtr(), C.CString(message))
	return (*ChatMessage)(msg)
}
