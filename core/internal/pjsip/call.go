package pjsip

import (
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type Call struct {
	pjsua2.Call
}

type call struct {
	call    pjsua2.Call
	account pjsua2.Account
	sys     doorpix.System
}

func NewCall(account pjsua2.Account, callId int, system doorpix.System) *Call {
	impl := &call{
		sys:     system,
		account: account,
	}
	director := pjsua2.NewDirectorCall(impl, account, callId)
	impl.call = director

	return &Call{
		Call: director,
	}
}

func (c *call) OnCallState(param pjsua2.OnCallStateParam) {
	slog.Debug("call state changed", "state", c.call.GetInfo().GetState(), "call id", c.call.GetId())

	switch c.call.GetInfo().GetState() {
	case pjsua2.PJSIP_INV_STATE_DISCONNECTED:
		acc := c.account.DirectorInterface().(*account)
		delete(acc.calls, c.call.GetId())
		pjsua2.DeleteDirectorCall(c.call)
	}
}

func (c *call) OnCallMediaState(param pjsua2.OnCallMediaStateParam) {
	audioMedia := c.call.GetAudioMedia(-1)
	if audioMedia == nil {
		slog.Debug("no audio media", "call id", c.call.GetId())
		return
	}

	audioPlayback := pjsua2.EndpointInstance().AudDevManager().GetPlaybackDevMedia()
	if audioPlayback != nil {
		audioMedia.StartTransmit(audioPlayback)
	} else {
		slog.Error("failed to get playback device media", "call id", c.call.GetId())
	}

	audioRecording := pjsua2.EndpointInstance().AudDevManager().GetCaptureDevMedia()
	if audioRecording != nil {
		audioRecording.StartTransmit(audioMedia)
	} else {
		slog.Error("failed to get capture device media", "call id", c.call.GetId())
	}
}
