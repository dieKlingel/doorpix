package config

import (
	"gopkg.in/yaml.v3"
)

type Event = string

type EventCollection map[Event][]Action

const (
	StartupEvent             Event = "startup"
	ShutdownEvent            Event = "shutdown"
	CallStateChangeEvent     Event = "call-state-change"
	CallStateConnectEvent    Event = "call-state-connect"
	CallStateDisconnectEvent Event = "call-state-disconnect"
	DtmfActionEvent          Event = "dtmf-action"
	CallIncomingEvent        Event = "call-incoming"
	NewMessageEvent          Event = "new-message"
)

func (collection EventCollection) UnmarshalYAML(node *yaml.Node) error {
	rawActionNodesByEvent := map[Event][]yaml.Node{}
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
