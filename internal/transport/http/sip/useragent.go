package sip

import (
	"github.com/dieklingel/doorpix/internal/transport/sip"
)

type UserAgent interface {
	AccountInfo() *sip.AccountInfo
}
