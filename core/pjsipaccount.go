package core

import (
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/config"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type PJSIPAccount struct {
	pjsua2.Account
}

type PJSIPAccountOverrides struct {
	account pjsua2.Account
	emitter *EventEmitter
}

func NewPJSIPAccount(emitter *EventEmitter) *PJSIPAccount {
	overrides := &PJSIPAccountOverrides{
		emitter: emitter,
	}
	director := pjsua2.NewDirectorAccount(overrides)
	overrides.account = director

	return &PJSIPAccount{
		Account: director,
	}
}

func (a *PJSIPAccountOverrides) OnRegState(param pjsua2.OnRegStateParam) {
	slog.Info("registration state changed", "reason", param.GetReason(), "status", param.GetStatus())
}

func (a *PJSIPAccountOverrides) OnIncomingCall(param pjsua2.OnIncomingCallParam) {
	slog.Info("incoming call", "uri", param.GetCallId())
	a.emitter.On(config.CallIncomingEvent)
}

func (a *PJSIPAccountOverrides) OnInstantMessage(param pjsua2.OnInstantMessageParam) {
	slog.Info("instant message", "uri", param.GetFromUri(), "message", param.GetMsgBody())
	a.emitter.On(config.NewMessageEvent)
}

func DeletePJSIPAccount(account *PJSIPAccount) {
	pjsua2.DeleteDirectorAccount(account.Account)
}
