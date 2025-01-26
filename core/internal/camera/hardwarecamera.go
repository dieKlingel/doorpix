package camera

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/go-gst/go-gst/gst"
)

var cameras = make(map[string]*HardwareCamera)

type HardwareCamera struct {
	gstPipeline   *gst.Pipeline
	gstSrcElement *gst.Element
	gstTeeElement *gst.Element

	mutex         sync.Mutex
	cameraCounter int
}

func NewHardwareCamera(fullIdentifier string) (*HardwareCamera, error) {
	identifier := strings.Split(fullIdentifier, " ")[0]

	if existingHardwareCamera, ok := cameras[identifier]; ok {
		return existingHardwareCamera, nil
	}

	gstPipeline, err := gst.NewPipelineFromString(fmt.Sprintf("%s ! tee name=tee", fullIdentifier))
	if err != nil {
		return nil, err
	}

	gstSrcElements, err := gstPipeline.GetSourceElements()
	if err != nil {
		return nil, err
	}
	if len(gstSrcElements) != 1 {
		return nil, fmt.Errorf("expected 1 source element, got %d", len(gstSrcElements))
	}
	gstSrcElement := gstSrcElements[0]

	gstTeeElement, err := gstPipeline.GetElementByName("tee")
	if err != nil {
		return nil, err
	}

	/*if err := gstPipeline.AddMany(gstSrcElement, gstTeeElement); err != nil {
		return nil, err
	}
	if err := gstSrcElement.Link(gstTeeElement); err != nil {
		return nil, err
	}*/

	hwcamera := &HardwareCamera{
		gstPipeline:   gstPipeline,
		gstSrcElement: gstSrcElement,
		gstTeeElement: gstTeeElement,
	}
	cameras[identifier] = hwcamera
	success := gstPipeline.GetPipelineBus().AddWatch(hwcamera.onNewPipelineBusMessage)
	if !success {
		slog.Error("failed to add watch to pipeline bus")
	}

	return hwcamera, nil
}

func LookUpHardwareCamera(identifier string) (*HardwareCamera, bool) {
	cam, ok := cameras[identifier]
	return cam, ok
}

func (c *HardwareCamera) pause() {
	slog.Debug("pausing camera")

	currentState := c.gstPipeline.GetCurrentState()
	if currentState == gst.StatePlaying {
		err := c.gstPipeline.SetState(gst.StatePaused)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func (c *HardwareCamera) play() {
	slog.Debug("playing camera")

	err := c.gstPipeline.SetState(gst.StatePlaying)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (c *HardwareCamera) stop() {
	slog.Debug("stopping camera")

	err := c.gstPipeline.SetState(gst.StateNull)
	if err != nil {
		slog.Error(err.Error())
	}

	// todo cleanup

}

func (c *HardwareCamera) onNewPipelineBusMessage(msg *gst.Message) bool {
	slog.Debug("received message", "type", msg.Type().String())

	switch msg.Type() {
	case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
		err := c.gstPipeline.BlockSetState(gst.StateNull)
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
