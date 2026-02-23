package sip

import (
	"log/slog"

	"github.com/dieklingel/go-pjproject/pjsua2"
)

type Call struct {
	account  *Account
	delegate pjsua2.Call
}

func NewCall(acc *Account) *Call {
	call := &Call{
		account: acc,
	}
	osThread.invoke(func() {
		delegate := pjsua2.NewDirectorCall(call, acc.delegate)
		call.delegate = delegate
	})

	return call
}

func NewCallFromId(acc *Account, id int) *Call {
	call := &Call{
		account: acc,
	}
	osThread.invoke(func() {
		delegate := pjsua2.NewDirectorCall(call, acc.delegate, id)
		call.delegate = delegate
	})

	return call
}

func (c *Call) OnCallState(param pjsua2.OnCallStateParam) {
	info := c.delegate.GetInfo()
	state := info.GetState()
	slog.Info("call state changed", "callId", info.GetId(), "state", info.GetStateText())

	switch state {
	case pjsua2.PJSIP_INV_STATE_DISCONNECTED:
		delete(c.account.calls, info.GetId())
	}
}

func (c *Call) OnCallMediaState(arg2 pjsua2.OnCallMediaStateParam) {
	// TODO(koifresh): wiring correct media connection, this just connects the default audio in/out to every call
	callAudioMedia := c.delegate.GetAudioMedia(-1)
	if callAudioMedia == nil {
		slog.Warn("the call has no audo media", "callId", c.delegate.GetId())
		return
	}

	captureDevMedia := osThread.endpoint.AudDevManager().GetCaptureDevMedia()
	if captureDevMedia != nil {
		captureDevMedia.StartTransmit(callAudioMedia)
	}

	playbackDevMedia := osThread.endpoint.AudDevManager().GetPlaybackDevMedia()
	if playbackDevMedia != nil {
		callAudioMedia.StartTransmit(playbackDevMedia)
	}
}

func (c *Call) Info() *CallInfo {
	var callInfo CallInfo
	osThread.invoke(func() {
		info := c.delegate.GetInfo()
		callInfo = CallInfo{
			Id:        info.GetId(),
			RemoteUri: info.GetRemoteUri(),
			State:     info.GetStateText(),
		}
	})

	return &callInfo
}
