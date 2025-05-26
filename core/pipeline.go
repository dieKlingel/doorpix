package core

type Pipeline struct {
	steps []PipelineStep
}

func (pipeline *Pipeline) Run(ctx map[string]any) error {
	for _, step := range pipeline.steps {
		if err := step.Run(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (pipeline *Pipeline) IsEmpty() bool {
	return len(pipeline.steps) == 0
}
