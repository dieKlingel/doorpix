package core

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"runtime"

	"github.com/dieklingel/doorpix/core/internal/drivers/appvideo"
	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/pjsip"
	"github.com/dieklingel/doorpix/core/internal/service"
	"github.com/dieklingel/doorpix/core/pkg/pjsua2"
)

type SIPClientProps struct {
	Username    string
	Password    string
	Realm       string
	Server      string
	StunServers []string
	VideoDevice string
	Whitelist   []string
}

type SIPClient struct {
	props SIPClientProps

	accountConfig pjsua2.AccountConfig
	account       *pjsip.Account

	eventemitter *eventemitter.EventEmitter
	command      chan pjsip.Command
	ctx          service.Context
	ready        chan bool
}

func NewSIPClient(eventemitter *eventemitter.EventEmitter, props SIPClientProps) *SIPClient {
	return &SIPClient{
		props:        props,
		eventemitter: eventemitter,
		command:      make(chan pjsip.Command),
		ctx:          service.NewContext(context.Background()),
		ready:        make(chan bool),
	}
}

func (s *SIPClient) Start() error {
	s.ctx.Lock()
	go func() {
		defer s.ctx.Unlock()

		s.exec()
	}()

	<-s.ready

	return nil
}

func (s *SIPClient) Stop() {
	s.ctx.CancelAndWait()
}

func (s *SIPClient) Invite(uris []string) error {
	for _, uri := range uris {
		cmd := &pjsip.InviteCommand{
			Uri: uri,
		}
		s.command <- cmd
		err := <-cmd.Error()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SIPClient) SendMessage(uris []string, message string) error {
	cmd := &pjsip.SendInstantMessageCommand{
		Uri:     uris[0],
		Message: message,
	}

	s.command <- cmd
	err := <-cmd.Error()

	return err
}

func (s *SIPClient) Hangup() error {
	cmd := &pjsip.HangupCommand{}
	s.command <- cmd
	err := <-cmd.Error()

	return err
}

func (s *SIPClient) exec() {
	// make pjsip run in the same thread
	runtime.LockOSThread()

	// global setup
	pjEndpointConfig := pjsua2.NewEpConfig()
	pjEndpointConfig.GetLogConfig().SetLevel(2)

	userAgentConfig := pjEndpointConfig.GetUaConfig()
	userAgentConfig.SetUserAgent("DoorPiX")
	for _, server := range s.props.StunServers {
		userAgentConfig.GetStunServer().Add(server)
	}

	pjEndpoint := pjsua2.NewEndpoint()
	pjEndpoint.LibCreate()
	pjEndpoint.LibInit(pjEndpointConfig)
	pjEndpoint.LibStart()

	// transport
	transport := pjsua2.NewTransportConfig()
	//transport.SetPort(5067)
	if res := pjsua2.EndpointInstance().TransportCreate(pjsua2.PJSIP_TRANSPORT_TLS, transport); res != 0 {
		slog.Error("Error creating transport", "status", res)
	}

	// account config
	s.accountConfig = pjsua2.NewAccountConfig()
	s.accountConfig.GetSipConfig().GetProxies().Add(fmt.Sprintf("sip:%s;hide;transport=tls", s.props.Server))
	s.accountConfig.SetIdUri(fmt.Sprintf("sip:%s@%s", s.props.Username, s.props.Realm))
	s.accountConfig.GetRegConfig().SetRegistrarUri(fmt.Sprintf("sip:%s", s.props.Realm))

	// video
	appvideo.SetCameraDevice(s.props.VideoDevice)
	videoDeviceIndex := -1
	videoDevices := pjsua2.EndpointInstance().VidDevManager().EnumDev2()
	for i := 0; i < int(videoDevices.Size()); i++ {
		if videoDevices.Get(i).GetName() == appvideo.GetCameraDeviceName() {
			videoDeviceIndex = i
			break
		}
	}
	if videoDeviceIndex >= 0 {
		s.accountConfig.GetVideoConfig().SetDefaultCaptureDevice(videoDeviceIndex)
	}

	s.accountConfig.GetVideoConfig().SetAutoTransmitOutgoing(true)
	s.accountConfig.GetVideoConfig().SetAutoShowIncoming(false)
	s.accountConfig.GetMediaConfig().SetSrtpSecureSignaling(1)
	s.accountConfig.GetMediaConfig().SetSrtpUse(pjsua2.PJMEDIA_SRTP_MANDATORY)

	cred := pjsua2.NewAuthCredInfo("digest", "*", s.props.Username, 0, s.props.Password)
	s.accountConfig.GetSipConfig().GetAuthCreds().Add(cred)

	// acccount
	s.account = pjsip.NewAccount(pjsip.AccountProps{
		OnInstantMessage: s.onInstantMessage,
		OnRegState:       s.onRegState,
		OnIncomingCall:   s.onIncomingCall,
	})
	s.account.Create(s.accountConfig)

	s.ready <- true

	for {
		select {
		case command := <-s.command:
			s.processCommand(command)

		case <-s.ctx.Done():
			// cleanup
			// TODO: refactor to a proper way to clean up pjsip
			pjsua2.DeleteAccount(s.account.Account)
			pjsua2.DeleteAccountConfig(s.accountConfig)
			pjsua2.EndpointInstance().LibStopWorkerThreads()
			pjsua2.EndpointInstance().LibDestroy()
			s.account = nil
			s.accountConfig = nil
			return
		}
	}
}

// this method is always called from the pjsip thread
func (s *SIPClient) processCommand(command pjsip.Command) {
	switch cmd := command.(type) {

	case *pjsip.SendInstantMessageCommand:
		err := s.account.SendInstantMessage(cmd.Uri, cmd.Message)
		cmd.Error() <- err

	case *pjsip.InviteCommand:
		err := s.account.Invite(cmd.Uri)
		cmd.Error() <- err

	case *pjsip.HangupCommand:
		err := s.account.Hangup()
		cmd.Error() <- err

	default:
		cmd.Error() <- fmt.Errorf("unknown command type: %T", cmd)
	}
}

func (s *SIPClient) onInstantMessage(param pjsua2.OnInstantMessageParam) {
	from := param.GetFromUri()
	body := param.GetMsgBody()
	contentType := param.GetContentType()

	data := map[string]any{
		"From":        from,
		"Body":        body,
		"ContentType": contentType,
	}

	eventPath := fmt.Sprintf("events/sip-message/%s", from)
	s.eventemitter.Emit(eventPath, data)
}

func (s *SIPClient) onIncomingCall(call *pjsip.Call, _ pjsua2.OnIncomingCallParam) {
	remoteUri := call.GetInfo().GetRemoteUri()
	regex := regexp.MustCompile("^.*?<(sip:|)(.*?)>$") // https://regex101.com/r/8ci6jN/1
	matches := regex.FindStringSubmatch(remoteUri)
	if len(matches) >= 2 {
		remoteUri = matches[2]
	}

	isWhitelisted := false
	for _, uri := range s.props.Whitelist {
		if uri == remoteUri {
			isWhitelisted = true
			break
		}
	}

	if isWhitelisted {
		slog.Info("accept incomming call", "uri", remoteUri)
		call.Accept()
	} else {
		slog.Info("decline incomming call, because the uri is not whitelisted", "uri", remoteUri)
		call.Decline()
	}
}

func (s *SIPClient) onRegState(param pjsua2.OnRegStateParam) {
}
