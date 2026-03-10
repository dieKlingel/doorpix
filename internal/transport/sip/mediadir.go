package sip

import (
	"fmt"

	"github.com/dieklingel/go-pjproject/pjsua2"
)

type MediaDir pjsua2.Pjmedia_dir

func (dir MediaDir) String() string {
	switch pjsua2.Pjmedia_dir(dir) {
	case pjsua2.PJMEDIA_DIR_CAPTURE:
		return "Capture"
	case pjsua2.PJMEDIA_DIR_CAPTURE_PLAYBACK:
		return "CapturePlayback"
	case pjsua2.PJMEDIA_DIR_CAPTURE_RENDER:
		return "CaptureRender"
	case pjsua2.PJMEDIA_DIR_DECODING:
		return "Decoding"
	case pjsua2.PJMEDIA_DIR_ENCODING:
		return "Encoding"
	case pjsua2.PJMEDIA_DIR_ENCODING_DECODING:
		return "EncodingDecoding"
	case pjsua2.PJMEDIA_DIR_NONE:
		return "None"
	case pjsua2.PJMEDIA_DIR_PLAYBACK:
		return "Playback"
	case pjsua2.PJMEDIA_DIR_RENDER:
		return "Render"
	}

	return fmt.Sprintf("Unknown(%d)", dir)
}
