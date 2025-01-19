package core

import "github.com/dieklingel/doorpix/core/internal/config"

type SIPPhone struct {
	Config *config.Config
}

func (s *SIPPhone) HandleEvent(action config.Action, event *Event) {
	switch action.(type) {
	case config.InviteAction:
		// Do something
	}
}

func (s *SIPPhone) Setup(emitter *EventEmitter) {
	emitter.Handler(s)
}

func (s *SIPPhone) Cleanup() {

}

func (s *SIPPhone) Exec() {

}
