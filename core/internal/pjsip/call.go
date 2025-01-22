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
	slog.Info("call state changed", "state", c.call.GetInfo().GetState())

	switch c.call.GetInfo().GetState() {
	case pjsua2.PJSIP_INV_STATE_DISCONNECTED:
		acc := c.account.DirectorInterface().(*account)
		delete(acc.calls, c.call.GetId())
		pjsua2.DeleteDirectorCall(c.call)
	}
}

func (c *call) OnCallMediaState(param pjsua2.OnCallMediaStateParam) {}
