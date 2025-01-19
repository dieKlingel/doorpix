package linphone

// #include <linphone/linphonecore.h>
import "C"

type AccountParams C.LinphoneAccountParams

func (params *AccountParams) cPtr() *C.LinphoneAccountParams {
	return (*C.LinphoneAccountParams)(params)
}

func (params *AccountParams) SetServerAddress(serverAddr *Address) {
	C.linphone_account_params_set_server_address(params.cPtr(), serverAddr.cPtr())
}

func (params *AccountParams) SetIdentityAddress(identityAddr *Address) {
	C.linphone_account_params_set_identity_address(params.cPtr(), identityAddr.cPtr())
}

func (params *AccountParams) SetTransport(transport TransportType) {
	C.linphone_account_params_set_transport(params.cPtr(), C.LinphoneTransportType(transport))
}

func (params *AccountParams) EnableRegister(register bool) {
	C.linphone_account_params_enable_register(params.cPtr(), boolToCBool(register))
}
