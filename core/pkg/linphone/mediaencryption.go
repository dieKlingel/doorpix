package linphone

// #include <linphone/linphonecore.h>
import "C"

type MediaEncryption = C.LinphoneMediaEncryption

const (
	MediaEncryptionNone MediaEncryption = C.LinphoneMediaEncryptionNone
	MediaEncryptionSRTP MediaEncryption = C.LinphoneMediaEncryptionSRTP
)
