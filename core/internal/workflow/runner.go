package workflow

import (
	"fmt"
	"path"
)

type Runner struct {
	registry *Registry
	source   PipelineSourceFunc
}

type PipelineSourceFunc func() ([]*Pipeline, error)

func NewRunner(registry *Registry) *Runner {
	return &Runner{
		registry: registry,
	}
}

func (r *Runner) SetPipelineSourceFunc(source PipelineSourceFunc) {
	r.source = source
}

func (r *Runner) SetPipelineSource(pipelines []*Pipeline) {
	r.source = func() ([]*Pipeline, error) {
		return pipelines, nil
	}
}

func (r *Runner) FindPipelines(s fmt.Stringer) ([]*Pipeline, error) {
	if r.source == nil {
		return nil, ErrSourceIsNil
	}

	pipelines, err := r.source()
	if err != nil {
		return nil, err
	}

	result := make([]*Pipeline, 0)
	for _, pipeline := range pipelines {
		match, err := path.Match(pipeline.Trigger, s.String())

		if err != nil {
			return nil, err
		}

		if match {
			result = append(result, pipeline)
		}
	}

	return result, nil
}

func (r *Runner) Run(pipeline *Pipeline) error {
	if r.registry == nil {
		return ErrRegistryIsNil
	}

	if pipeline == nil {
		return ErrPipelineIsNil
	}

	delegates := make([]StepDelegate, 0, len(pipeline.Steps))

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
