package workflowtest

import "github.com/dieklingel/doorpix/core/internal/workflow"

type NilProvider struct {
	ParseError error
	RunError   error
}

func (p *NilProvider) Parse(step workflow.Step) (workflow.StepDelegate[any], error) {
	if p.ParseError != nil {
		return workflow.StepDelegate[any]{}, p.ParseError
	}

	return workflow.StepDelegate[any]{
		Step: step,
		Execute: func(ctx *workflow.Context) error {
			return p.RunError
		},
	}, nil
}
