package sip

import (
	"fmt"
	"log/slog"

	cameradriver "github.com/dieklingel/doorpix/internal/transport/sip/driver/camera"
	"github.com/dieklingel/go-pjproject/pjsua2"
)

type Account struct {
	delegate pjsua2.Account
	calls    map[int]*Call
}

type AccountProps struct {
	Username string
	Password string
	Realm    string
	Domain   string
}

func NewAccount(props AccountProps) (*Account, error) {
	acc := &Account{
		calls: make(map[int]*Call),
	}
	username := props.Username
	realm := props.Realm
	domain := props.Domain
	password := props.Password

	if len(domain) == 0 {
		domain = realm
	}

	var err error = nil
	osThread.invoke(func() {
		delegate := pjsua2.NewDirectorAccount(acc)
		acc.delegate = delegate

		transport := pjsua2.NewTransportConfig()
		if res := pjsua2.EndpointInstance().TransportCreate(pjsua2.PJSIP_TRANSPORT_UDP, transport); res != 0 {
			err = fmt.Errorf("could not create transport of type tls")
			return
		}

		cfg := pjsua2.NewAccountConfig()
		cfg.GetSipConfig().GetProxies().Add(fmt.Sprintf("sip:%s;hide;transport=udp", domain))
		cfg.SetIdUri(fmt.Sprintf("sip:%s@%s", username, realm))
		cfg.GetRegConfig().SetRegistrarUri(fmt.Sprintf("sip:%s", domain))

		cfg.GetVideoConfig().SetDefaultCaptureDevice(-1)
		cfg.GetVideoConfig().SetAutoTransmitOutgoing(true)
		cfg.GetVideoConfig().SetAutoShowIncoming(false)
		cfg.GetMediaConfig().SetSrtpSecureSignaling(0)
		cfg.GetMediaConfig().SetSrtpUse(pjsua2.PJMEDIA_SRTP_OPTIONAL)

		cdevs := osThread.endpoint.VidDevManager().EnumDev2()
		for i := range cdevs.Size() {
			dev := cdevs.Get(int(i))
			if dev.GetName() == cameradriver.DeviceName() {
				cfg.GetVideoConfig().SetDefaultCaptureDevice(dev.GetId())
			}
			slog.Info("account: camera capture device found", "name", dev.GetName(), "id", dev.GetId())
		}

		cred := pjsua2.NewAuthCredInfo("digest", "*", username, 0, password)
		cfg.GetSipConfig().GetAuthCreds().Add(cred)

		acc.delegate.Create(cfg)
	})

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (acc *Account) OnIncomingCall(param pjsua2.OnIncomingCallParam) {
	id := param.GetCallId()
	slog.Info("received incomming call", "callId", id)

	call := NewCallFromId(acc, id)
	acc.calls[id] = call

	op := pjsua2.NewCallOpParam()
	op.SetStatusCode(pjsua2.PJSIP_SC_OK)

	call.delegate.Answer(op)
}

func (acc *Account) OnRegStarted(arg2 pjsua2.OnRegStartedParam) {

}

func (acc *Account) OnRegState(arg2 pjsua2.OnRegStateParam) {

}

func (acc *Account) OnIncomingSubscribe(arg2 pjsua2.OnIncomingSubscribeParam) {

}

func (acc *Account) OnInstantMessage(arg2 pjsua2.OnInstantMessageParam) {

}

func (acc *Account) OnInstantMessageStatus(arg2 pjsua2.OnInstantMessageStatusParam) {

}

func (acc *Account) OnSendRequest(arg2 pjsua2.OnSendRequestParam) {

}

func (acc *Account) OnTypingIndication(arg2 pjsua2.OnTypingIndicationParam) {

}

func (acc *Account) OnMwiInfo(arg2 pjsua2.OnMwiInfoParam) {

}
