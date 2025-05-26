package actions

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"text/template"

	"gopkg.in/yaml.v3"
)

type ConditionAction struct {
	Condition template.Template
	Then      []Action
	Else      []Action
}

func (a *ConditionAction) Execute(data map[string]any) (bool, error) {
	condition := bytes.Buffer{}
	err := a.Condition.Execute(&condition, data)
	if err != nil {
		return false, err
	}

	command := fmt.Sprintf("%s && echo -n true || echo -n false", condition.String())
	value, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return false, err
	}

	conditionMet, err := strconv.ParseBool(string(value))
	if err != nil {
		return false, err
	}

	if conditionMet {
		return true, nil
	}

	return false, nil
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
