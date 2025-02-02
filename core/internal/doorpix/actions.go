package doorpix

import (
	"fmt"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

type Action interface{}

type ConditionAction struct {
	Condition template.Template
	Then      []Action
	Else      []Action
}

func (a *ConditionAction) UnmarshalYAML(node *yaml.Node) error {
	raw := struct {
		Condition YamlStringTemplate `yaml:"condition"`
		Then      []yaml.Node        `yaml:"then"`
		Else      []yaml.Node        `yaml:"else"`
	}{}

	err := node.Decode(&raw)
	if err != nil {
		return err
	}

	a.Condition = (template.Template)(raw.Condition)

	a.Then = make([]Action, len(raw.Then))
	for i, thenAction := range raw.Then {
		action, err := newActionFromNode(thenAction)
		if err != nil {
			return err
		}
		a.Then[i] = action
	}

	a.Else = make([]Action, len(raw.Else))
	for i, elseAction := range raw.Else {
		action, err := newActionFromNode(elseAction)
		if err != nil {
			return err
		}
		a.Else[i] = action
	}

	return nil
}

type SleepAction struct {
	Duration time.Duration `yaml:"sleep"`
}

func (a *SleepAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		Sleep string `yaml:"sleep"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	duration, err := time.ParseDuration(action.Sleep)
	if err != nil {
		return err
	}

	a.Duration = duration
	return nil
}

type LogAction struct {
	Message *template.Template
}

func (a *LogAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		Log string `yaml:"log"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	tmpl, err := template.New("log").Parse(action.Log)
	if err != nil {
		return err
	}

	a.Message = tmpl
	return nil
}

type EvalAction struct {
	Expressions []*template.Template `yaml:"eval"`
}

func (a *EvalAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		Expressions YamlScalarOrList[YamlStringTemplate] `yaml:"eval"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	a.Expressions = make([]*template.Template, len(action.Expressions))
	for i, expr := range action.Expressions {
		a.Expressions[i] = (*template.Template)(&expr)
	}
	return nil
}

type InviteAction struct {
	UriTemplates []template.Template
}

func (a *InviteAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		UriTemplates YamlScalarOrList[YamlStringTemplate] `yaml:"invite"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	a.UriTemplates = make([]template.Template, len(action.UriTemplates))
	for i := range action.UriTemplates {
		a.UriTemplates[i] = (template.Template)(action.UriTemplates[i])
	}
	return nil
}

type MessageAction struct {
	UriTemplates    []template.Template `yaml:"to"`
	MessageTemplate template.Template   `yaml:"message"`
}

func (a *MessageAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		UriTemplates YamlScalarOrList[YamlStringTemplate] `yaml:"to"`
		Message      YamlStringTemplate                   `yaml:"message"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	a.UriTemplates = make([]template.Template, len(action.UriTemplates))
	for i, number := range action.UriTemplates {
		a.UriTemplates[i] = (template.Template)(number)
	}
	a.MessageTemplate = (template.Template)(action.Message)
	return nil
}

type HangupAction struct{}

type RPCAction struct {
	Timeout  time.Duration
	Spec     map[string]string
	Multiple bool
}

func (a *RPCAction) UnmarshalYAML(node *yaml.Node) error {
	type rpc struct {
		Timeout  string            `yaml:"timeout"`
		Spec     map[string]string `yaml:"spec"`
		Multiple bool              `yaml:"multiple"`
	}

	action := struct {
		Rpc rpc `yaml:"rpc"`
	}{
		Rpc: rpc{
			Timeout: "5s",
		},
	}

	if err := node.Decode(&action); err != nil {
		return err
	}

	timeout, err := time.ParseDuration(action.Rpc.Timeout)
	if err != nil {
		return err
	}
	a.Timeout = timeout
	a.Spec = action.Rpc.Spec
	a.Multiple = action.Rpc.Multiple
	return nil
}

type PublishAction struct {
	Topic   template.Template
	Message template.Template
}

func (a *PublishAction) UnmarshalYAML(node *yaml.Node) error {
	action := struct {
		Topic   YamlStringTemplate `yaml:"topic"`
		Message YamlStringTemplate `yaml:"message"`
	}{}

	err := node.Decode(&action)
	if err != nil {
		return err
	}

	a.Topic = (template.Template)(action.Topic)
	a.Message = (template.Template)(action.Message)
	return nil
}

func newActionFromNode(node yaml.Node) (Action, error) {
	if node.Kind == yaml.MappingNode {
		raw := map[string]any{}
		if err := node.Decode(&raw); err != nil {
			return nil, err
		}

		if raw["condition"] != nil {
			action := ConditionAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["sleep"] != nil {
			action := SleepAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["log"] != nil {
			action := LogAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["hangup"] != nil {
			action := HangupAction{}
			err := node.Decode(&action)
			return HangupAction{}, err
		}
		if raw["invite"] != nil {
			action := InviteAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["eval"] != nil {
			action := EvalAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["message"] != nil && raw["to"] != nil {
			action := MessageAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["topic"] != nil && raw["message"] != nil {
			action := PublishAction{}
			err := node.Decode(&action)
			return action, err
		}
		if raw["rpc"] != nil {
			action := RPCAction{}
			err := node.Decode(&action)
			return action, err
		}
	} else if node.Kind == yaml.ScalarNode {
		if node.Value == "hangup" {
			return HangupAction{}, nil
		}
	}

	return nil, fmt.Errorf("could not infer action type in line %d", node.Line)
}
