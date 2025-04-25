package core

import (
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/pjsip"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type PJSIPService struct {
	Config doorpix.Config

	account       *pjsip.Account
	accountConfig pjsua2.AccountConfig
	pjaction      chan pjAction
}

type pjAction struct {
}

func (service *PJSIPService) Name() string {
	return "pjsip-service"
}

/*
func (service *PJSIPService) Run(action doorpix.Action, event *doorpix.ActionHook) {
	service.pjaction <- pjAction{action, event}
}

func (service *PJSIPService) init() {
	config := pjsua2.NewEpConfig()
	config.GetLogConfig().SetLevel(2)

	ua := config.GetUaConfig()
	ua.SetUserAgent("DoorPiX")

	for _, server := range service.Config.SIPPhone.StunServers {
		ua.GetStunServer().Add(server)
	}

	endpoint := pjsua2.NewEndpoint()
	endpoint.LibCreate()
	endpoint.LibInit(config)
	endpoint.LibStart()

	// needs to be called after pjlib is started
	// TODO: refactor to a proper way to initialize the video device
	appvideo.SetCameraDevice(service.Config.Camera.Device)
}

func (service *PJSIPService) deinit() {
	slog.Debug("shutting down pjsip service")

	// TODO: refactor to a proper way to clean up pjsip
	pjsua2.DeleteAccount(service.account.Account)
	pjsua2.DeleteAccountConfig(service.accountConfig)
	pjsua2.EndpointInstance().LibStopWorkerThreads()
	pjsua2.EndpointInstance().LibDestroy()
	service.account = nil
	service.accountConfig = nil

	slog.Debug("successfully shut down pjsip service")
}

func (service *PJSIPService) run(action doorpix.Action, event *doorpix.ActionHook) {
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

func (service *PJSIPService) Exec(ctx context.Context, wg *sync.WaitGroup) error {
	slog.Debug("exec pjsip service")

	wg.Add(1)
	go func() {
		defer wg.Done()
		runtime.LockOSThread()

		service.exec(ctx)
	}()

	return nil
}

func (service *PJSIPService) exec(ctx context.Context) {
	service.init()

	// transport
	transport := pjsua2.NewTransportConfig()
	//transport.SetPort(5067)
	if res := pjsua2.EndpointInstance().TransportCreate(pjsua2.PJSIP_TRANSPORT_TLS, transport); res != 0 {
		slog.Error("Error creating transport")
	}

	// account config
	service.accountConfig = pjsua2.NewAccountConfig()
	service.accountConfig.GetSipConfig().GetProxies().Add(fmt.Sprintf("sip:%s;hide;transport=tls", service.Config.SIPPhone.Server))
	service.accountConfig.SetIdUri(fmt.Sprintf("sip:%s@%s", service.Config.SIPPhone.Username, service.Config.SIPPhone.Realm))
	service.accountConfig.GetRegConfig().SetRegistrarUri(fmt.Sprintf("sip:%s", service.Config.SIPPhone.Realm))

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

	cred := pjsua2.NewAuthCredInfo("digest", "*", service.Config.SIPPhone.Username, 0, service.Config.SIPPhone.Password)
	service.accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	// acccount
	service.account = pjsip.NewAccount(service.Config, service.Emit)
	service.account.Create(service.accountConfig)

	for {
		select {
		case <-ctx.Done():
			service.deinit()
			return
		case action := <-service.pjaction:
			service.run(action.action, action.event)
		}
	}
}
*/
