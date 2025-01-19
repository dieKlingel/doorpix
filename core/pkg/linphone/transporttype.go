package linphone

// #include <linphone/linphonecore.h>
import "C"
import (
	"log/slog"
	"strings"
)

type TransportType C.LinphoneTransportType

const (
	TransportTypeUdp  TransportType = C.LinphoneTransportUdp
	TransportTypeTcp  TransportType = C.LinphoneTransportTcp
	TransportTypeTls  TransportType = C.LinphoneTransportTls
	TransportTypeDtls TransportType = C.LinphoneTransportDtls
)

func (t *TransportType) cPtr() *C.LinphoneTransportType {
	return (*C.LinphoneTransportType)(t)
}

func TransportTypeFromString(transport string) TransportType {
	switch strings.ToLower(transport) {
	case "udp":
		return TransportTypeUdp
	case "tcp":
		return TransportTypeTcp
	case "tls":
		return TransportTypeTls
	case "dtls":
		return TransportTypeDtls
	default:
		slog.Warn("the transport '%s' could not be parsed, udp will be used instead", "transport", transport)
		return TransportTypeUdp
	}
}
