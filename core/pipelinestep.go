package core

type PipelineStep struct {
	delegate func(ctx map[string]any) error
}

func (step *PipelineStep) Run(ctx map[string]any) error {
	if step.delegate == nil {
		return nil
	}

	return step.delegate(ctx)
}
