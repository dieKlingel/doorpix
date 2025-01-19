package linphone

type CoreUserData struct {
	onGlobalStateChanged       OnGlobalStateChangedCallback
	onRegistrationStateChanged OnRegistrationStateChangedCallback
	onCallStateChanged         OnCallStateChangedCallback
	onDtmfReceived             OnDtmfReceivedCallback
}
