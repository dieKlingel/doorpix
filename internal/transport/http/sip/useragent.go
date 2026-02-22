package sip

import (
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

type UserAgent interface {
	Calls() []sip.CallInfo
	AccountInfo() *sip.AccountInfo
}
