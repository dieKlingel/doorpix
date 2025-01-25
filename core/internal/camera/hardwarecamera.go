package camera

import (
	"log/slog"
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

func NewHardwareCamera(identifier string) (*HardwareCamera, error) {
	if existingHardwareCamera, ok := cameras[identifier]; ok {
		return existingHardwareCamera, nil
	}

	gstPipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	gstSrcElement, err := gst.NewElement(identifier)
	if err != nil {
		return nil, err
	}

	gstTeeElement, err := gst.NewElement("tee")
	if err != nil {
		return nil, err
	}

	if err := gstPipeline.AddMany(gstSrcElement, gstTeeElement); err != nil {
		return nil, err
	}
	if err := gstSrcElement.Link(gstTeeElement); err != nil {
		return nil, err
	}

	hwcamera := &HardwareCamera{
		gstPipeline:   gstPipeline,
		gstSrcElement: gstSrcElement,
		gstTeeElement: gstTeeElement,
	}
	cameras[identifier] = hwcamera
	gstPipeline.GetPipelineBus().AddWatch(hwcamera.onNewPipelineBusMessage)

	return hwcamera, nil
}

func LookUpHardwareCamera(identifier string) (*HardwareCamera, bool) {
	cam, ok := cameras[identifier]
	return cam, ok
}

func (c *HardwareCamera) pause() {
	currentState := c.gstPipeline.GetCurrentState()
	if currentState == gst.StatePlaying {
		err := c.gstPipeline.SetState(gst.StatePaused)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func (c *HardwareCamera) play() {
	err := c.gstPipeline.SetState(gst.StatePlaying)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (c *HardwareCamera) stop() {
	err := c.gstPipeline.SetState(gst.StateNull)
	if err != nil {
		slog.Error(err.Error())
	}

	// todo cleanup

}

func (c *HardwareCamera) onNewPipelineBusMessage(msg *gst.Message) bool {
	switch msg.Type() {
	case gst.MessageEOS: // When end-of-stream is received flush the pipeling and stop the main loop
		c.gstPipeline.BlockSetState(gst.StateNull)
	case gst.MessageError: // Error messages are always fatal
		err := msg.ParseError()

		slog.Error(err.Error())
		if debug := err.DebugString(); debug != "" {
			slog.Debug(debug)
		}
	}
	return true
}
