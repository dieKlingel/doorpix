package pjsip

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/dieklingel/pjproject-dist/go/pjsua2"
)

// Account is a wrapper around the pjsua2.Account struct. It provides a user
// friendly interface to interact with the pjsip account tailored for the
// doorpix system
type Account struct {
	pjsua2.Account
}

type AccountProps struct {
	OnInstantMessage func(param pjsua2.OnInstantMessageParam)
	OnRegState       func(param pjsua2.OnRegStateParam)
	OnIncomingCall   func(call *Call, param pjsua2.OnIncomingCallParam)
}

// Invite sends an invite to the given uri. The uri must be a valid sip uri
// (e.g. sip:example@example.org or example@example.org)
func (a *Account) Invite(uri string) error {
	if !strings.HasPrefix(uri, "sip:") {
		uri = fmt.Sprintf("sip:%s", uri)
	}

	if pjsua2.EndpointInstance().UtilVerifySipUri(uri) != 0 {
		return fmt.Errorf("%s is not a valid sip uri", uri)
	}

	account := a.Account.DirectorInterface().(*account)

	param := pjsua2.NewCallOpParam()
	call := NewCall(a)
	call.MakeCall(uri, param)

	account.calls[call.GetId()] = call
	return nil
}

// Hangup hangs all active calls
func (a *Account) Hangup() error {
	account := a.Account.DirectorInterface().(*account)
	for index, call := range account.calls {
		if call.IsActive() {
			call.Hangup(pjsua2.NewCallOpParam())
			delete(account.calls, index)
		}
	}

	return nil
}

// SendInstantMessage sends an instant message to the given uri. The uri must be
// a valid sip uri (e.g. sip:example@example.org or example@example.org)
func (a *Account) SendInstantMessage(uri string, message string) error {
	if !strings.HasPrefix(uri, "sip:") {
		uri = fmt.Sprintf("sip:%s", uri)
	}

	if pjsua2.EndpointInstance().UtilVerifySipUri(uri) != 0 {
		return fmt.Errorf("%s is not valid sip uri", uri)
	}

	buddy := pjsua2.NewBuddy()
	config := pjsua2.NewBuddyConfig()

	config.SetUri(uri)
	buddy.Create(a, config)

	messageParam := pjsua2.NewSendInstantMessageParam()
	messageParam.SetContent(message)
	buddy.SendInstantMessage(messageParam)
	return nil
}

type account struct {
	account pjsua2.Account
	calls   map[int]*Call
	props   AccountProps
}

func NewAccount(props AccountProps) *Account {
	impl := &account{
		props: props,
		calls: make(map[int]*Call),
	}
	director := pjsua2.NewDirectorAccount(impl)
	impl.account = director

	return &Account{
		Account: director,
	}
}

func (a *account) OnRegState(param pjsua2.OnRegStateParam) {
	slog.Info("registration state changed", "reason", param.GetReason(), "status", param.GetStatus())

	a.props.OnRegState(param)
}

func (a *account) OnIncomingCall(param pjsua2.OnIncomingCallParam) {
	call := NewCallWithId(a.account, param.GetCallId())

	slog.Debug("incoming call", "from", call.GetInfo().GetRemoteUri(), "callId", param.GetCallId())
	a.calls[param.GetCallId()] = call

	a.props.OnIncomingCall(call, param)
}

func (a *account) OnInstantMessage(param pjsua2.OnInstantMessageParam) {
	slog.Debug("new message received", "from", param.GetFromUri(), "type", param.GetContentType(), "message", param.GetMsgBody())

	a.props.OnInstantMessage(param)
}

func DeletePJSIPAccount(account *Account) {
	pjsua2.DeleteDirectorAccount(account.Account)
}
