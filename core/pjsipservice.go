package core

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/drivers/appvideo"
	"github.com/dieklingel/doorpix/core/internal/pjsip"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type PJSIPService struct {
	System doorpix.System

	account       *pjsip.Account
	accountConfig pjsua2.AccountConfig
}

func (service *PJSIPService) HandleEvent(action doorpix.Action, event *doorpix.Event) {
	if !pjsua2.EndpointInstance().LibIsThreadRegistered() {
		pjsua2.EndpointInstance().LibRegisterThread("")
	}

	switch action := action.(type) {
	case doorpix.InviteAction:
		for _, uriTemplate := range action.UriTemplates {
			var uri bytes.Buffer
			if err := uriTemplate.Execute(&uri, event.Data); err != nil {
				slog.Error(err.Error())
				continue
			}

			err := service.account.Invite(uri.String())
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

			err := service.account.SendInstantMessage(uri.String(), message.String())
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}
}

func (service *PJSIPService) Setup() {
	service.System.Bus.Handler(service)

	// init
	service.init()
}

func (service *PJSIPService) init() {
	config := pjsua2.NewEpConfig()
	config.GetLogConfig().SetLevel(2)

	ua := config.GetUaConfig()
	ua.SetUserAgent("DoorPiX")

	for _, server := range service.System.Config.SIPPhone.StunServers {
		ua.GetStunServer().Add(server)
	}

	endpoint := pjsua2.NewEndpoint()
	endpoint.LibCreate()
	endpoint.LibInit(config)
	endpoint.LibStart()

	// needs to be called after pjlib is started
	// TODO: refactor to a proper way to initialize the video device
	appvideo.SetCameraDevice(service.System.Config.Camera.Device)
}

func (service *PJSIPService) Exec() {
	if !pjsua2.EndpointInstance().LibIsThreadRegistered() {
		pjsua2.EndpointInstance().LibRegisterThread("exec")
		slog.Debug("register thread for pjsip", "name", "exec", "for", "PJSIPPhone")
	}

	// transport
	transport := pjsua2.NewTransportConfig()
	//transport.SetPort(5067)
	if res := pjsua2.EndpointInstance().TransportCreate(pjsua2.PJSIP_TRANSPORT_TLS, transport); res != 0 {
		slog.Error("Error creating transport")
	}

	// account config
	service.accountConfig = pjsua2.NewAccountConfig()
	service.accountConfig.GetSipConfig().GetProxies().Add(fmt.Sprintf("sip:%s;hide;transport=tls", service.System.Config.SIPPhone.Server))
	service.accountConfig.SetIdUri(fmt.Sprintf("sip:%s@%s", service.System.Config.SIPPhone.Username, service.System.Config.SIPPhone.Realm))
	service.accountConfig.GetRegConfig().SetRegistrarUri(fmt.Sprintf("sip:%s", service.System.Config.SIPPhone.Realm))

	videoDeviceIndex := -1
	videoDevices := pjsua2.EndpointInstance().VidDevManager().EnumDev2()
	for i := 0; i < int(videoDevices.Size()); i++ {
		if videoDevices.Get(i).GetName() == appvideo.GetCameraDeviceName() {
			videoDeviceIndex = i
			break
		}
	}
	if videoDeviceIndex >= 0 {
		service.accountConfig.GetVideoConfig().SetDefaultCaptureDevice(videoDeviceIndex)
	}

	service.accountConfig.GetVideoConfig().SetAutoTransmitOutgoing(true)
	service.accountConfig.GetVideoConfig().SetAutoShowIncoming(false)
	service.accountConfig.GetMediaConfig().SetSrtpSecureSignaling(1)
	service.accountConfig.GetMediaConfig().SetSrtpUse(pjsua2.PJMEDIA_SRTP_MANDATORY)

	cred := pjsua2.NewAuthCredInfo("digest", "*", service.System.Config.SIPPhone.Username, 0, service.System.Config.SIPPhone.Password)
	service.accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	// acccount
	service.account = pjsip.NewAccount(service.System)
	service.account.Create(service.accountConfig)
}

func (service *PJSIPService) Cleanup() {

}
