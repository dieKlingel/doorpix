package doorpix

import (
	"gopkg.in/yaml.v3"
)

type EventType = string

type EventCollection map[EventType][]Action

const (
	StartupEvent             EventType = "startup"
	ShutdownEvent            EventType = "shutdown"
	CallStateChangeEvent     EventType = "call-state-change"
	CallStateConnectEvent    EventType = "call-state-connect"
	CallStateDisconnectEvent EventType = "call-state-disconnect"
	DtmfActionEvent          EventType = "dtmf-action"
	CallIncomingEvent        EventType = "call-incoming"
	NewMessageEvent          EventType = "new-message"
)

func (collection EventCollection) UnmarshalYAML(node *yaml.Node) error {
	rawActionNodesByEvent := map[EventType][]yaml.Node{}
	if err := node.Decode(&rawActionNodesByEvent); err != nil {
		return err
	}

	for event, rawActionNodes := range rawActionNodesByEvent {
		actions := make([]Action, 0, len(rawActionNodes))

		for _, rawActionNode := range rawActionNodes {
			action, err := newActionFromNode(rawActionNode)
			if err != nil {
				return err
			}
			actions = append(actions, action)
		}

		collection[event] = actions
	}

	return nil
}
