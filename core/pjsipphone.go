package core

import (
	"bytes"
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

func (p *PJSIPPhone) HandleEvent(action doorpix.Action, event *doorpix.Event) {
	if !pjsua2.EndpointInstance().LibIsThreadRegistered() {
		pjsua2.EndpointInstance().LibRegisterThread("pjsipphone")
	}

	switch action := action.(type) {
	case doorpix.InviteAction:
	case doorpix.MessageAction:
		var msg bytes.Buffer
		if err := action.Message.Execute(&msg, event.Data); err != nil {
			slog.Error(err.Error())
			break
		}

		for _, number := range action.Numbers {
			buddy := pjsua2.NewBuddy()
			config := pjsua2.NewBuddyConfig()

			var uri bytes.Buffer
			if err := number.Execute(&uri, event.Data); err != nil {
				slog.Error(err.Error())
				continue
			}
			slog.Info("Sending message to", "number", uri.String())

			config.SetUri(fmt.Sprintf("sip:%s", uri.String()))
			buddy.Create(p.account, config)

			message := pjsua2.NewSendInstantMessageParam()
			message.SetContent(msg.String())
			buddy.SendInstantMessage(message)
		}
	}
}

func (p *PJSIPPhone) Setup() {
	p.System.Bus.Handler(p)

	// global config
	config := pjsua2.NewEpConfig()
	config.GetLogConfig().SetLevel(2)

	ua := config.GetUaConfig()
	ua.SetUserAgent("DoorPiX")

	ua.GetStunServer().Add("stun.l.google.com:19302")
	ua.GetStunServer().Add("stun.linphone.org:3478")

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

	p.accountConfig.GetVideoConfig().SetAutoTransmitOutgoing(true)
	p.accountConfig.GetVideoConfig().SetAutoShowIncoming(false)
	p.accountConfig.GetMediaConfig().SetSrtpSecureSignaling(1)
	p.accountConfig.GetMediaConfig().SetSrtpUse(pjsua2.PJMEDIA_SRTP_MANDATORY)

	cred := pjsua2.NewAuthCredInfo("digest", "*", p.System.Config.SIPPhone.Username, 0, p.System.Config.SIPPhone.Password)
	p.accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	// acccount
	p.account = pjsip.NewAccount(p.System)
	p.account.Create(p.accountConfig)
}

func (p *PJSIPPhone) Cleanup() {

}
