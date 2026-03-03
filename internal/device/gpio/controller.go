package gpio

import (
	"context"
	"fmt"

	"github.com/dieklingel/doorpix/internal/oplog"
	"github.com/warthog618/go-gpiocdev"
)

type ControllerProps struct {
	Chip   string
	Inputs []int
}

type Controller struct {
	chip   string
	inputs []int
	lines  *gpiocdev.Lines
}

func NewController(props ControllerProps) *Controller {
	return &Controller{
		chip:   props.Chip,
		inputs: props.Inputs,
	}
}

func (c *Controller) Run() error {
	lines, err := gpiocdev.RequestLines(
		c.chip,
		c.inputs,
		gpiocdev.AsInput,
		gpiocdev.WithBothEdges,
		gpiocdev.WithEventHandler(c.OnGpioEvent),
	)

	if err != nil {
		return err
	}

	c.lines = lines
	return nil
}

func (c *Controller) OnGpioEvent(evt gpiocdev.LineEvent) {
	edge := "unknown"
	switch evt.Type {
	case gpiocdev.LineEventFallingEdge:
		edge = "falling"
	case gpiocdev.LineEventRisingEdge:
		edge = "rising"
	}

	oplog.Dispatch(
		fmt.Sprintf("system/doorpix/gpio/%d/%s", evt.Offset, edge),
		"edge", edge,
		"input", evt.Offset,
		"chip", c.chip,
		"timestamp", evt.Timestamp,
	)
}

func (c *Controller) Stop(ctx context.Context) error {
	if c.lines != nil {
		err := c.lines.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
