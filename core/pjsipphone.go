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
		for _, uriTemplate := range action.UriTemplates {
			var uri bytes.Buffer
			if err := uriTemplate.Execute(&uri, event.Data); err != nil {
				slog.Error(err.Error())
				continue
			}

			err := p.account.Invite(uri.String())
			if err != nil {
				slog.Error(err.Error())
			}
		}

	case doorpix.MessageAction:
		var message bytes.Buffer
		if err := action.MessageTemplate.Execute(&message, event.Data); err != nil {
			slog.Error(err.Error())
			break
		}

		for _, uriTemplate := range action.UriTemplates {
			var uri bytes.Buffer
			if err := uriTemplate.Execute(&uri, event.Data); err != nil {
				slog.Error(err.Error())
				continue
			}

			err := p.account.SendInstantMessage(uri.String(), message.String())
			if err != nil {
				slog.Error(err.Error())
			}
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

	videoDeviceIndex := -1
	videoDevices := pjsua2.EndpointInstance().VidDevManager().EnumDev2()
	for i := 0; i < int(videoDevices.Size()); i++ {
		if videoDevices.Get(i).GetName() == "DoorPiX Emulated Video Device" {
			videoDeviceIndex = i
			break
		}
	}
	if videoDeviceIndex >= 0 {
		slog.Info("Setting video device", "name", videoDevices.Get(videoDeviceIndex).GetName())
		p.accountConfig.GetVideoConfig().SetDefaultCaptureDevice(videoDeviceIndex)
	}

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
