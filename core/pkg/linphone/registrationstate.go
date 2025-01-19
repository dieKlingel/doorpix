package linphone

// #include <linphone/linphonecore.h>
import "C"

type LinphoneRegistrationState C.LinphoneRegistrationState

const (
	LinphoneRegistrationStateNone     LinphoneRegistrationState = C.LinphoneRegistrationNone
	LinphoneRegistrationStateOk       LinphoneRegistrationState = C.LinphoneRegistrationOk
	LinphoneRegistrationStateFailed   LinphoneRegistrationState = C.LinphoneRegistrationFailed
	LinphoneRegistrationStateProgress LinphoneRegistrationState = C.LinphoneRegistrationProgress
	LinphoneRegistrationStateCleared  LinphoneRegistrationState = C.LinphoneRegistrationCleared
)
