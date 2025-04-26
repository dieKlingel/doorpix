package core

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dieklingel/doorpix/core/internal/eventemitter"
	"github.com/dieklingel/doorpix/core/internal/service"
	"github.com/warthog618/go-gpiocdev"
)

type GPIOControllerProps struct {
	Pins map[int]string
	Chip string
}

type GPIOController struct {
	props        GPIOControllerProps
	eventemitter *eventemitter.EventEmitter

	chip *gpiocdev.Chip
	ctx  service.Context
}

func NewGPIOController(eventemitter *eventemitter.EventEmitter, props GPIOControllerProps) *GPIOController {
	return &GPIOController{
		props:        props,
		eventemitter: eventemitter,

		ctx: service.NewContext(context.Background()),
	}
}

func (g *GPIOController) Start() error {
	chip, err := gpiocdev.NewChip(g.props.Chip)
	if err != nil {
		return err
	}

	g.chip = chip
	g.exec()

	return nil
}

func (g *GPIOController) Stop() error {
	g.ctx.CancelAndWait()

	return nil
}

func (g *GPIOController) exec() {
	g.ctx.Lock()
	go func() {
		defer g.ctx.Unlock()

		lines, err := g.chip.RequestLines([]int{0}, gpiocdev.AsInput, gpiocdev.WithEventHandler(g.onLineChange))
		if err != nil {
			slog.Warn("failed to request GPIO lines", "error", err)
		}
		defer lines.Close()

		<-g.ctx.Done()
	}()
}

func (g *GPIOController) onLineChange(event gpiocdev.LineEvent) {
	slog.Debug("GPIO event", "event", event)

	pinNumber := event.Offset
	pinName := g.props.Pins[pinNumber]
	var eventType string

	switch event.Type {
	case gpiocdev.LineEventRisingEdge:
		eventType = "rising"
	case gpiocdev.LineEventFallingEdge:
		eventType = "falling"

	default:
		slog.Warn("unknown GPIO event type", "type", event.Type)
		return
	}

	eventData := map[string]any{
		"Edge":      eventType,
		"Pin":       pinNumber,
		"Name":      pinName,
		"Timestamp": event.Timestamp,
	}
	eventPath := fmt.Sprintf("events/gpio/%s/%s", pinName, eventType)
	g.eventemitter.Emit(eventPath, eventData)
}
