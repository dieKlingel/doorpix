package sip

import (
	"fmt"

	"github.com/dieklingel/go-pjproject/pjsua2"
)

type MediaType pjsua2.Pjmedia_type

func (t MediaType) String() string {
	switch pjsua2.Pjmedia_type(t) {
	case pjsua2.PJMEDIA_TYPE_APPLICATION:
		return "Application"
	case pjsua2.PJMEDIA_TYPE_AUDIO:
		return "Audio"
	case pjsua2.PJMEDIA_TYPE_NONE:
		return "None"
	case pjsua2.PJMEDIA_TYPE_UNKNOWN:
		return "Unknown"
	case pjsua2.PJMEDIA_TYPE_VIDEO:
		return "Video"
	}

	return fmt.Sprintf("Unknown(%d)", t)
}
