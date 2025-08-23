package workflowtest

import "github.com/dieklingel/doorpix/core/internal/workflow"

type NilProvider struct {
	ParseError error
	RunError   error
}

func (p *NilProvider) Parse(step workflow.Step) (workflow.StepDelegate, error) {
	if p.ParseError != nil {
		return workflow.StepDelegate{}, p.ParseError
	}

	return workflow.StepDelegate{
		Step: step,
		Execute: func(ctx *workflow.Context) error {
			return p.RunError
		},
	}, nil
}
