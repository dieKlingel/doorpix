package doorpix

import (
	"github.com/dieklingel/doorpix/core/internal/actions"
	"gopkg.in/yaml.v3"
)

type EventType = string

type EventCollection map[EventType][]actions.Action

const (
	StartupEvent             EventType = "startup"
	ShutdownEvent            EventType = "shutdown"
	CallStateChangeEvent     EventType = "call-state-change"
	CallStateConnectEvent    EventType = "call-state-connect"
	CallStateDisconnectEvent EventType = "call-state-disconnect"
	DtmfActionEvent          EventType = "dtmf-action"
	CallIncomingEvent        EventType = "call-incoming"
	MqttMessageEvent         EventType = "mqtt-message"
	NewMessageEvent          EventType = "new-message"
	APIRingEvent             EventType = "api:ring"
	APIUnlockEvent           EventType = "api:unlock"
)

func (collection EventCollection) UnmarshalYAML(node *yaml.Node) error {
	rawActionNodesByEvent := map[EventType][]yaml.Node{}
	if err := node.Decode(&rawActionNodesByEvent); err != nil {
		return err
	}

	for event, rawActionNodes := range rawActionNodesByEvent {
		parsedActions := make([]actions.Action, 0, len(rawActionNodes))

		for _, rawActionNode := range rawActionNodes {
			action, err := actions.Parse(rawActionNode)
			if err != nil {
				return err
			}
			parsedActions = append(parsedActions, action)
		}

		collection[event] = parsedActions
	}

	return nil
}
