package linphone

// #include <linphone/linphonecore.h>
import "C"

type LoggingService C.LinphoneLoggingService

func (service *LoggingService) Ptr() *LoggingService {
	return service
}

func (service *LoggingService) cPtr() *C.LinphoneLoggingService {
	return (*C.LinphoneLoggingService)(service)
}
