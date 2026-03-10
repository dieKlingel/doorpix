package sip

import (
	"fmt"

	"github.com/dieklingel/go-pjproject/pjsua2"
)

type CallMediaStatus pjsua2.Pjsua_call_media_status

func (s CallMediaStatus) String() string {
	switch pjsua2.Pjsua_call_media_status(s) {
	case pjsua2.PJSUA_CALL_MEDIA_ACTIVE:
		return "Active"
	case pjsua2.PJSUA_CALL_MEDIA_ERROR:
		return "Error"
	case pjsua2.PJSUA_CALL_MEDIA_LOCAL_HOLD:
		return "LocalHold"
	case pjsua2.PJSUA_CALL_MEDIA_NONE:
		return "None"
	case pjsua2.PJSUA_CALL_MEDIA_REMOTE_HOLD:
		return "RemoteHold"
	}

	return fmt.Sprintf("Unknown(%d)", s)
}
