package core

import (
	"fmt"
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/config"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type PJSIPPhone struct {
	Config  *config.Config
	emitter *EventEmitter

	// pjsua2
	endpoint      pjsua2.Endpoint
	account       *PJSIPAccount
	accountConfig pjsua2.AccountConfig
}

func (p *PJSIPPhone) HandleEvent(config config.Action, event *Event) {

}

func (p *PJSIPPhone) Setup(emitter *EventEmitter) {
	emitter.Handler(p)
	p.emitter = emitter

	// global config
	config := pjsua2.NewEpConfig()
	config.GetLogConfig().SetLevel(2)
	p.endpoint = pjsua2.NewEndpoint()
	p.endpoint.LibCreate()
	p.endpoint.LibInit(config)

}

func (p *PJSIPPhone) Exec() {
	if !pjsua2.EndpointInstance().LibIsThreadRegistered() {
		slog.Info("Registering thread")
		pjsua2.EndpointInstance().LibRegisterThread("exec")
	}

	// transport
	transport := pjsua2.NewTransportConfig()
	//transport.SetPort(5067)
	if res := p.endpoint.TransportCreate(pjsua2.PJSIP_TRANSPORT_TLS, transport); res != 0 {
		slog.Error("Error creating transport")
	}

	// account config
	p.accountConfig = pjsua2.NewAccountConfig()
	p.accountConfig.GetSipConfig().GetProxies().Add(fmt.Sprintf("sip:%s;hide;transport=tls", p.Config.SIPPhone.Server))
	p.accountConfig.SetIdUri(fmt.Sprintf("sip:%s@%s", p.Config.SIPPhone.Username, p.Config.SIPPhone.Realm))
	p.accountConfig.GetRegConfig().SetRegistrarUri(fmt.Sprintf("sip:%s", p.Config.SIPPhone.Realm))

	cred := pjsua2.NewAuthCredInfo("digest", "*", p.Config.SIPPhone.Username, 0, p.Config.SIPPhone.Password)
	p.accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	// acccount
	p.account = NewPJSIPAccount(p.emitter)
	p.account.Create(p.accountConfig)
}

func (p *PJSIPPhone) Cleanup() {

}
