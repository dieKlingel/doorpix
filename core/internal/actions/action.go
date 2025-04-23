package actions

import "gopkg.in/yaml.v3"

type Action interface{}

func Parse(node yaml.Node) (Action, error) {
	action, err := newActionFromNode(node)

	if err != nil {
		return nil, err
	}

	return action, nil
}
