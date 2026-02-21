package sip

import (
	"context"

	"github.com/dieklingel/go-pjproject/v2/pjsua2"
)

type UserAgentProps struct {
	Username string
	Password string
	Realm    string
	Domain   string
}

type UserAgent struct {
	props   UserAgentProps
	account *Account
}

// Only one UserAgent should be created at all
func NewUserAgent(props UserAgentProps) UserAgent {
	return UserAgent{
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
