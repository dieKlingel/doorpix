package doorpi

import "github.com/dieklingel/doorpix/internal/transport/sip"

type UserAgent interface {
	Invite(uri string) (*sip.CallInfo, error)
}
