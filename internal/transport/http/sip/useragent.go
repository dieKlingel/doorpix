package sip

import (
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

type UserAgent interface {
	Calls() []sip.CallInfo
	CallById(id int) *sip.CallInfo
	AccountInfo() *sip.AccountInfo
	Invite(uri string) (*sip.CallInfo, error)
	Hangup(id int)
}
