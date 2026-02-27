package sip

import (
	"context"
	"errors"

	cameradriver "github.com/dieklingel/doorpix/internal/transport/sip/driver/camera"
	"github.com/dieklingel/go-pjproject/pjsua2"
)

type UserAgentProps struct {
	Username string
	Password string
	Realm    string
	Domain   string
	Webcam   cameradriver.Webcam
}

type UserAgent struct {
	props   UserAgentProps
	account *Account
}

// Only one UserAgent should be created at all
func NewUserAgent(props UserAgentProps) *UserAgent {
	return &UserAgent{
		props: props,
	}
}

func (ua *UserAgent) Serve() error {
	go osThread.run()

	success := false
	osThread.invoke(func() {
		success = osThread.endpoint.LibIsThreadRegistered()
	})

	if !success {
		return NativeThreadError
	}

	var err error
	osThread.invoke(func() {
		err = cameradriver.Register("Internal", ua.props.Webcam)
	})
	if err != nil {
		return err
	}

	acc, err := NewAccount(AccountProps{
		Username: ua.props.Username,
		Password: ua.props.Password,
		Realm:    ua.props.Realm,
		Domain:   ua.props.Domain,
	})

	ua.account = acc
	return err
}

func (ua *UserAgent) AccountInfo() *AccountInfo {
	if ua.account == nil {
		return nil
	}

	var accountInfo pjsua2.AccountInfo
	osThread.invoke(func() {
		accountInfo = ua.account.delegate.GetInfo()
	})

	return &AccountInfo{
		Uri:        accountInfo.GetUri(),
		IsActive:   accountInfo.GetRegIsActive(),
		StatusText: accountInfo.GetRegStatusText(),
	}
}

func (ua *UserAgent) Calls() []CallInfo {
	if ua.account == nil {
		return make([]CallInfo, 0)
	}

	var callInfos []CallInfo = make([]CallInfo, 0, len(ua.account.calls))
	for _, call := range ua.account.calls {
		info := call.Info()
		callInfos = append(callInfos, *info)
	}

	return callInfos
}

func (ua *UserAgent) Invite(uri string) (*CallInfo, error) {
	if ua.account == nil {
		return nil, ErrNotReady
	}

	var err error
	osThread.invoke(func() {
		if osThread.endpoint.UtilVerifySipUri(uri) != 0 {
			err = errors.Join(ErrInvalidUri, errors.New(uri))
		}
	})
	if err != nil {
		return nil, err
	}

	call := NewCall(ua.account)

	osThread.invoke(func() {
		op := pjsua2.NewCallOpParam()
		call.delegate.MakeCall(uri, op)

		id := call.delegate.GetId()
		ua.account.calls[id] = call
	})

	info := call.Info()
	return info, nil
}

func (ua *UserAgent) CallById(id int) *CallInfo {
	call, exists := ua.account.calls[id]
	if !exists {
		return nil
	}

	return call.Info()
}

func (ua *UserAgent) Hangup(id int) {
	call, exists := ua.account.calls[id]
	if !exists {
		return
	}

	delete(call.account.calls, id)

	osThread.invoke(func() {
		op := pjsua2.NewCallOpParam()
		call.delegate.Hangup(op)
	})
}

func (ua *UserAgent) Shutdown(ctx context.Context) error {
	finished := make(chan struct{})

	go osThread.invoke(func() {
		osThread.endpoint.LibDestroy()
		finished <- struct{}{}
		close(finished)
	})

	select {
	case <-finished:
		osThread.done <- struct{}{}
		return nil
	case <-ctx.Done():
		return UserAgentShutdownError
	}
}
