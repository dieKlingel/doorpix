package camera

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type GstDriver struct {
	gstPipeline *gst.Pipeline

	mutex   sync.Mutex
	counter int
	valves  map[string]bool
}

func NewGstDriver(input string) (*GstDriver, error) {
	gst.Init(nil)

	gstPipeline, err := gst.NewPipelineFromString(input)
	if err != nil {
		return nil, err
	}

	driver := &GstDriver{
		gstPipeline: gstPipeline,
		valves:      make(map[string]bool),
	}

	success := gstPipeline.GetPipelineBus().AddWatch(driver.onNewPipelineBusMessage)
	if !success {
		return nil, fmt.Errorf("failed to add bus watch to pipeline")
	}

	return driver, nil
}

func (d *GstDriver) GetAppSinkByName(name string) (*app.Sink, error) {
	rawName := fmt.Sprintf("appsink-%s", strings.ToLower(name))
	element, err := d.gstPipeline.GetElementByName(rawName)
	if err != nil {
		return nil, err
	}

	appSink := app.SinkFromElement(element)
	if appSink == nil {
		return nil, fmt.Errorf("cannot convert the element %s to appsink", name)
	}

	return appSink, nil
}

func (d *GstDriver) Start(name string) error {
	rawName := fmt.Sprintf("valve-%s", strings.ToLower(name))
	element, err := d.gstPipeline.GetElementByName(rawName)
	if err != nil {
		return err
	}

	drop, err := element.GetProperty("drop")
	if err != nil {
		return err
	}

	if drop := drop.(bool); drop {
		err = element.SetProperty("drop", false)
		if err != nil {
			return err
		}
	}

	d.valves[name] = false
	return d.setState(gst.StatePlaying)
}

func (d *GstDriver) Stop(name string) error {
	rawName := fmt.Sprintf("valve-%s", strings.ToLower(name))
	element, err := d.gstPipeline.GetElementByName(rawName)
	if err != nil {
		return err
	}

	err = element.SetProperty("drop", true)
	if err != nil {
		return err
	}

	d.valves[name] = true

	if !d.hasOpenValve() {
		return d.setState(gst.StateNull)
	}

	return nil
}

func (d *GstDriver) hasOpenValve() bool {
	for _, drop := range d.valves {
		if !drop {
			return true
		}
	}

	return false
}

func (d *GstDriver) setState(state gst.State) error {
	currentState := d.gstPipeline.GetCurrentState()

	if currentState == state {
		slog.Debug("ignoring request to change camera pipeline state, it is already in the desired state", "state", state.String())
		return nil
	}

	slog.Debug("changing camera pipeline state", "from", currentState.String(), "to", state)

	err := d.gstPipeline.BlockSetState(state)
	if err != nil {
		return err
	}

	return nil
}

func (d *GstDriver) onNewPipelineBusMessage(msg *gst.Message) bool {
	slog.Debug("received new gst pipeline bus message", "type", msg.Type().String())

	switch msg.Type() {
	case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
		err := d.gstPipeline.BlockSetState(gst.StateNull)
		if err != nil {
			slog.Error(err.Error())
		}
	case gst.MessageError: // Error messages are always fatal
		err := msg.ParseError()

		slog.Error(err.Error())
		if debug := err.DebugString(); debug != "" {
			slog.Debug(debug)
		}
	}

	return true
}
