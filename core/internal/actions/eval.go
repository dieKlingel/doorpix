package actions

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type EvalAction struct {
	Expressions []*template.Template `yaml:"eval"`
}

func (action *EvalAction) Execute(data map[string]any) error {
	builder := strings.Builder{}

	for _, expression := range action.Expressions {
		buffer := bytes.Buffer{}
		err := expression.Execute(&buffer, data)
		if err != nil {
			return err
		}

		builder.WriteString(buffer.String())
		builder.WriteString("\n")
	}

	command := builder.String()
	output, err := exec.Command("sh", "-c", command).CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Print(string(output))
	return nil
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
