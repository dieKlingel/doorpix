package core

import (
	"fmt"
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/pjsip"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type PJSIPPhone struct {
	System doorpix.System

	// pjsua2
	endpoint      pjsua2.Endpoint
	account       *pjsip.Account
	accountConfig pjsua2.AccountConfig
}

func (p *PJSIPPhone) HandleEvent(config doorpix.Action, event *doorpix.Event) {

}

func (p *PJSIPPhone) Setup() {
	p.System.Bus.Handler(p)

	// global config
	config := pjsua2.NewEpConfig()
	config.GetLogConfig().SetLevel(2)
	p.endpoint = pjsua2.NewEndpoint()
	p.endpoint.LibCreate()
	p.endpoint.LibInit(config)
	p.endpoint.LibStart()
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
	p.accountConfig.GetSipConfig().GetProxies().Add(fmt.Sprintf("sip:%s;hide;transport=tls", p.System.Config.SIPPhone.Server))
	p.accountConfig.SetIdUri(fmt.Sprintf("sip:%s@%s", p.System.Config.SIPPhone.Username, p.System.Config.SIPPhone.Realm))
	p.accountConfig.GetRegConfig().SetRegistrarUri(fmt.Sprintf("sip:%s", p.System.Config.SIPPhone.Realm))

	cred := pjsua2.NewAuthCredInfo("digest", "*", p.System.Config.SIPPhone.Username, 0, p.System.Config.SIPPhone.Password)
	p.accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	// acccount
	p.account = pjsip.NewAccount(p.System)
	p.account.Create(p.accountConfig)
}

func (p *PJSIPPhone) Cleanup() {

}
