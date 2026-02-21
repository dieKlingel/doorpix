package sip

import "errors"

var NativeThreadError = errors.New("could not invoke native thread")
var UserAgentShutdownError = errors.New("the user agent could not be shut down successfully")
