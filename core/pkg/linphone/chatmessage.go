package linphone

// #include <linphone/linphonecore.h>
import "C"

type ChatMessage C.LinphoneChatMessage

func (m *ChatMessage) cPtr() *C.LinphoneChatMessage {
	return (*C.LinphoneChatMessage)(m)
}

func (m *ChatMessage) Send() {
	C.linphone_chat_message_send(m.cPtr())
}
