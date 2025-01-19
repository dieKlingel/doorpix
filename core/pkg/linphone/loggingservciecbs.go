package linphone

// #include <linphone/linphonecore.h>
import "C"

type LoggingServiceCbs C.LinphoneLoggingServiceCbs

func (c *LoggingServiceCbs) Ptr() *LoggingServiceCbs {
	return (*LoggingServiceCbs)(c)
}

func (cbs *LoggingServiceCbs) cPtr() *C.LinphoneLoggingServiceCbs {
	return (*C.LinphoneLoggingServiceCbs)(cbs)
}
