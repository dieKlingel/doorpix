package workflow

import "fmt"

type Runner struct {
	registry *Registry
}

func NewRunner(registry *Registry) *Runner {
	return &Runner{
		registry: registry,
	}
}

func (r *Runner) Run(pipeline *Pipeline) error {
	if r.registry == nil {
		return ErrRegistryIsNil
	}

	if pipeline == nil {
		return ErrPipelineIsNil
	}

	var delegates []StepDelegate = make([]StepDelegate, 0, len(pipeline.Steps))

	for _, step := range pipeline.Steps {
		provider, exists := r.registry.GetProvider(step.Uses)
		if !exists {
			return fmt.Errorf("provider not found for step: %s; %w", step.Uses, ErrProviderNotFound)
		}

		stepDelegate, err := provider.Parse(step)
		if err != nil {
			return fmt.Errorf("error parsing step %s: %w", step.Uses, err)
		}

		delegates = append(delegates, stepDelegate)
	}

	ctx := NewContext()
	for _, delegate := range delegates {
		err := delegate.Execute(ctx)
		if err != nil {
			return fmt.Errorf("error running step %s: %w", delegate.Step.Uses, err)
		}
	}

	return nil
}
