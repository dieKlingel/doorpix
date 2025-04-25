package pjsip

import (
	"fmt"
	"strings"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

// Account is a wrapper around the pjsua2.Account struct. It provides a user
// friendly interface to interact with the pjsip account tailored for the
// doorpix system
type Account struct {
	pjsua2.Account
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

	//åaccount := a.Account.DirectorInterface().(*account)

	//param := pjsua2.NewCallOpParam()
	//call := NewCall(a, -1, account.config)
	//call.MakeCall(uri, param)

	//account.calls[call.GetId()] = call
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
	config  doorpix.Config
}

/*

func NewAccount(config doorpix.Config, emit doorpix.Emit) *Account {
	impl := &account{
		config: config,
		emit:   emit,
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
	call := NewCall(a.account, param.GetCallId(), a.config, a.emit)
	remoteUri := call.GetInfo().GetRemoteUri()
	regex := regexp.MustCompile("^\".*?\"\\s<sip:(.*?)>$")
	matches := regex.FindStringSubmatch(remoteUri)
	if len(matches) >= 1 {
		remoteUri = matches[1]
	}

	slog.Info("incoming call received", "uri", remoteUri)
	a.emit(doorpix.CallIncomingEvent, nil)

	callParam := pjsua2.NewCallOpParam()
	callParam.SetStatusCode(pjsua2.PJSIP_SC_DECLINE)
	a.calls[param.GetCallId()] = call

	for _, uri := range a.config.SIPPhone.Whitelist {
		if uri == remoteUri {
			callParam.SetStatusCode(pjsua2.PJSIP_SC_OK)
			break
		}
	}

	if callParam.GetStatusCode() == pjsua2.PJSIP_SC_DECLINE {
		slog.Info("decline incomming call, because the uri is not whitelisted", "uri", remoteUri)
	} else {
		slog.Info("accept incomming call", "uri", remoteUri)
	}

	call.Answer(callParam)
}

func (a *account) OnInstantMessage(param pjsua2.OnInstantMessageParam) {
	slog.Debug("new message received", "from", param.GetFromUri(), "type", param.GetContentType())

	a.emit(doorpix.NewMessageEvent, nil)
}

func DeletePJSIPAccount(account *Account) {
	pjsua2.DeleteDirectorAccount(account.Account)
}
*/
