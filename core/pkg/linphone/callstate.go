package linphone

// #include <linphone/linphonecore.h>
import "C"

type CallState C.LinphoneCallState

const (
	CallStateIdle                 CallState = C.LinphoneCallIdle
	CallStateIncomingReceived     CallState = C.LinphoneCallIncomingReceived
	CallStatePushIncomingReceived CallState = C.LinphoneCallPushIncomingReceived
	CallStateOutgoingInit         CallState = C.LinphoneCallOutgoingInit
	CallStateOutgoingProgress     CallState = C.LinphoneCallOutgoingProgress
	CallStateOutgoingRinging      CallState = C.LinphoneCallOutgoingRinging
	CallStateOutgoingEarlyMedia   CallState = C.LinphoneCallOutgoingEarlyMedia
	CallStateConnected            CallState = C.LinphoneCallConnected
	CallStateStreamsRunning       CallState = C.LinphoneCallStreamsRunning
	CallStatePausing              CallState = C.LinphoneCallPausing
	CallStatePaused               CallState = C.LinphoneCallPaused
	CallStateResuming             CallState = C.LinphoneCallResuming
	CallStateRefered              CallState = C.LinphoneCallRefered
	CallStateError                CallState = C.LinphoneCallError
	CallStateEnd                  CallState = C.LinphoneCallEnd
	CallStatePausedByRemote       CallState = C.LinphoneCallPausedByRemote
	CallStateUpdatedByRemote      CallState = C.LinphoneCallUpdatedByRemote
	CallStateIncomingEarlyMedia   CallState = C.LinphoneCallIncomingEarlyMedia
	CallStateUpdating             CallState = C.LinphoneCallUpdating
	CallStateReleased             CallState = C.LinphoneCallReleased
	CallStateEarlyUpdatedByRemote CallState = C.LinphoneCallEarlyUpdatedByRemote
	CallStateEarlyUpdating        CallState = C.LinphoneCallEarlyUpdating
)
