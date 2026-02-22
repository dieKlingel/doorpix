package sip

import (
	"log/slog"

	"github.com/dieklingel/go-pjproject/pjsua2"
)

type Call struct {
	account  *Account
	delegate pjsua2.Call
}

type CallProps struct {
	Account *Account
	Id      int
}

func NewCall(props CallProps) *Call {
	call := &Call{
		account: props.Account,
	}
	osThread.invoke(func() {
		delegate := pjsua2.NewDirectorCall(call, props.Account.delegate, props.Id)
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
