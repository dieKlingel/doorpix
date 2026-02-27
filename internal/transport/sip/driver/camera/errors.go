package camera

import "errors"

var ErrAlreadyInitialized = errors.New("sip camera driver: already initalized")
var ErrRegistrationFailed = errors.New("sip camera driver: camera driver registration failed")
