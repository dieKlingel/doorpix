package pjsip

import (
	"log/slog"
	"regexp"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type Account struct {
	pjsua2.Account
}

type account struct {
	account pjsua2.Account
	system  doorpix.System
	calls   map[int]*Call
}

func NewAccount(system doorpix.System) *Account {
	impl := &account{
		system: system,
		calls:  make(map[int]*Call),
	}
	director := pjsua2.NewDirectorAccount(impl)
	impl.account = director

	return &Account{
		Account: director,
	}
}

func (a *account) OnRegState(param pjsua2.OnRegStateParam) {
	slog.Info("registration state changed", "reason", param.GetReason(), "status", param.GetStatus())
}

func (a *account) OnIncomingCall(param pjsua2.OnIncomingCallParam) {
	call := NewCall(a.account, param.GetCallId(), a.system)
	remoteUri := call.GetInfo().GetRemoteUri()
	regex := regexp.MustCompile("^\".*?\"\\s<sip:(.*?)>$")
	matches := regex.FindStringSubmatch(remoteUri)
	if len(matches) >= 1 {
		remoteUri = matches[1]
	}

	slog.Info("incoming call", "uri", remoteUri)
	a.system.Bus.On(doorpix.CallIncomingEvent)

	callParam := pjsua2.NewCallOpParam()
	callParam.SetStatusCode(pjsua2.PJSIP_SC_DECLINE)
	a.calls[param.GetCallId()] = call

	for _, uri := range a.system.Config.SIPPhone.Whitelist {
		if uri == remoteUri {
			callParam.SetStatusCode(pjsua2.PJSIP_SC_OK)
			break
		}
	}

	call.GetInfo().GetSetting().SetAudioCount(1)
	call.GetInfo().GetSetting().SetVideoCount(1)
	call.Answer(callParam)
}

func (a *account) OnInstantMessage(param pjsua2.OnInstantMessageParam) {
	slog.Info("instant message", "uri", param.GetFromUri(), "message", param.GetMsgBody())
	a.system.Bus.On(doorpix.NewMessageEvent)
}

func DeletePJSIPAccount(account *Account) {
	pjsua2.DeleteDirectorAccount(account.Account)
}
