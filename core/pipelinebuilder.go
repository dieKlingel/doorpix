package core

import (
	"fmt"
	"os"
	"reflect"

	"github.com/dieklingel/doorpix/core/internal/actions"
)

type MqttMessagePublisher interface {
	Publish(topic string, message string) error
}

type DependencyError struct {
	Action actions.Action
}

func (e *DependencyError) Error() string {
	return fmt.Sprintf("no handler for %v is configured", reflect.TypeOf(e.Action))
}

type UnknownActionError struct {
	Action actions.Action
}

func (e *UnknownActionError) Error() string {
	return fmt.Sprintf("unknown action type: %v", reflect.TypeOf(e.Action))
}

type PipelineBuilder struct {
	workflow             []actions.Action
	softphone            Softphone
	mqttMessagePublisher MqttMessagePublisher
}

func NewPipelineBuilder() *PipelineBuilder {
	pipeline := &PipelineBuilder{}
	return pipeline
}

func (builder *PipelineBuilder) WithSoftphone(softphone Softphone) *PipelineBuilder {
	builder.softphone = softphone
	return builder
}

func (builder *PipelineBuilder) WithMqttMessagePublisher(publisher MqttMessagePublisher) *PipelineBuilder {
	builder.mqttMessagePublisher = publisher
	return builder
}

func (builder *PipelineBuilder) WithActions(workflow []actions.Action) *PipelineBuilder {
	builder.workflow = workflow
	return builder
}

func (builder *PipelineBuilder) Clone() *PipelineBuilder {
	clone := &PipelineBuilder{
		workflow:             make([]actions.Action, len(builder.workflow)),
		softphone:            builder.softphone,
		mqttMessagePublisher: builder.mqttMessagePublisher,
	}
	copy(clone.workflow, builder.workflow)
	return clone
}

func (builder *PipelineBuilder) Build() (*Pipeline, error) {
	var steps = make([]PipelineStep, 0, len(builder.workflow))

	for _, action := range builder.workflow {
		switch action := action.(type) {
		case actions.ConditionAction:
			step, err := builder.newConditionStep(action)
			if err != nil {
				return nil, err
			}

			steps = append(steps, *step)
		case actions.LogAction:
			step := builder.newLogStep(action)
			steps = append(steps, *step)
		case actions.PublishAction:
			step, err := builder.newPublishStep(action)
			if err != nil {
				return nil, err
			}
			steps = append(steps, *step)
		case actions.MessageAction:
			step, err := builder.newMessageStep(action)
			if err != nil {
				return nil, err
			}

			steps = append(steps, *step)
		default:
			return nil, &UnknownActionError{Action: action}
		}
	}

	return &Pipeline{
		steps: steps,
	}, nil
}

func (builder *PipelineBuilder) newConditionStep(action actions.ConditionAction) (*PipelineStep, error) {
	thenPipeline, err := builder.Clone().WithActions(action.Then).Build()
	if err != nil {
		return nil, err
	}
	elsePipeline, err := builder.Clone().WithActions(action.Else).Build()
	if err != nil {
		return nil, err
	}

	step := PipelineStep{
		delegate: func(ctx map[string]any) error {
			conditionMet, err := action.Execute(ctx)
			if err != nil {
				return err
			}

			if conditionMet {
				return thenPipeline.Run(ctx)
			} else {
				return elsePipeline.Run(ctx)
			}
		},
	}

	return &step, nil
}

func (builder *PipelineBuilder) newLogStep(action actions.LogAction) *PipelineStep {
	return &PipelineStep{
		delegate: func(ctx map[string]any) error {
			return action.Execute(os.Stdout, ctx)
		},
	}
}

func (builder *PipelineBuilder) newPublishStep(action actions.PublishAction) (*PipelineStep, error) {
	if IsNil(builder.mqttMessagePublisher) {
		return nil, &DependencyError{Action: action}
	}

	return &PipelineStep{
		delegate: func(ctx map[string]any) error {
			return action.Execute(builder.mqttMessagePublisher, ctx)
		},
	}, nil
}

func (builder *PipelineBuilder) newMessageStep(action actions.MessageAction) (*PipelineStep, error) {
	if IsNil(builder.softphone) {
		return nil, &DependencyError{Action: action}
	}

	return &PipelineStep{
		delegate: func(ctx map[string]any) error {
			return action.Execute(builder.softphone, ctx)
		},
	}, nil
}

func IsNil(i any) bool {
	if i == nil {
		return true
	}

	return reflect.ValueOf(i).IsNil()
}
