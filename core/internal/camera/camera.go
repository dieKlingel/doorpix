package camera

import (
	"fmt"
	"log/slog"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

func init() {
	gst.Init(nil)
}

type Camera struct {
	gstQueueElement       *gst.Element
	gstAdditionalElements []*gst.Element
	gstAppSink            *app.Sink
	gstTeeSrcPad          *gst.Pad
	gstQueueSinkPad       *gst.Pad

	hardwareCamera *HardwareCamera
	frameChannel   chan []byte
	isRunning      bool
}

func NewFromString(device string, elements ...*gst.Element) (*Camera, error) {
	hardwareCamera, err := NewHardwareCamera(device)
	if err != nil {
		return nil, err
	}

	return NewFromHardwareCamera(hardwareCamera, elements...)
}

func NewFromHardwareCamera(hardwareCamera *HardwareCamera, elements ...*gst.Element) (*Camera, error) {
	queue, err := gst.NewElement("queue")
	if err != nil {
		return nil, err
	}

	appsink, err := app.NewAppSink()
	if err != nil {
		return nil, err
	}

	queueSinkPad := queue.GetStaticPad("sink")
	if queueSinkPad == nil {
		return nil, fmt.Errorf("queue has no static pad")
	}

	teeSrcPad := hardwareCamera.gstTeeElement.GetRequestPad("src_%u")
	if teeSrcPad == nil {
		return nil, fmt.Errorf("tee has no request pad")
	}

	camera := &Camera{
		gstQueueElement:       queue,
		gstAdditionalElements: elements,
		gstAppSink:            appsink,
		gstTeeSrcPad:          teeSrcPad,
		gstQueueSinkPad:       queueSinkPad,
		hardwareCamera:        hardwareCamera,
		isRunning:             false,
	}

	appsink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: camera.onNewSample,
	})

	/*gstTempPipeline, err := gst.NewPipelineFromString(pipeline)
	if err != nil {
		return nil, err
	}*/

	//gstElements, err := gstTempPipeline.GetElementsSorted()

	return camera, nil
}

func (c *Camera) SetProperty(key string, value any) error {
	err := c.hardwareCamera.gstSrcElement.SetProperty(key, value)
	return err
}

func (c *Camera) Start() error {
	slog.Debug("starting camera", "camera", c)

	c.hardwareCamera.mutex.Lock()
	slog.Debug("successfully locked the hardware device", "camera", c, "driver", c.hardwareCamera)

	defer func() {
		c.hardwareCamera.mutex.Unlock()
		slog.Debug("successfully unlocked the hardware device", "camera", c, "driver", c.hardwareCamera)
	}()

	if c.isRunning {
		slog.Debug("camera already running", "camera", c)
		return nil
	}

	c.isRunning = true
	c.hardwareCamera.cameraCounter++
	slog.Debug("increment the hardware device counter", "camera", c, "driver", c.hardwareCamera, "counter", c.hardwareCamera.cameraCounter)

	c.hardwareCamera.pause()

	gstElementsSize := 2 + len(c.gstAdditionalElements)
	allGstElements := make([]*gst.Element, gstElementsSize)
	allGstElements[0] = c.gstQueueElement
	for i, element := range c.gstAdditionalElements {
		allGstElements[i+1] = element
	}
	allGstElements[gstElementsSize-1] = c.gstAppSink.Element

	err := c.hardwareCamera.gstPipeline.AddMany(allGstElements...)
	if err != nil {
		return err
	}

	err = gst.ElementLinkMany(allGstElements...)
	if err != nil {
		return err
	}

	linkResult := c.gstTeeSrcPad.Link(c.gstQueueSinkPad)
	if linkResult != gst.PadLinkOK {
		return fmt.Errorf("could not link tee and queue")
	}

	c.hardwareCamera.play()
	c.frameChannel = make(chan []byte)
	return nil
}

func (c *Camera) Stop() error {
	slog.Debug("stopping camera", "camera", c)

	c.hardwareCamera.mutex.Lock()
	slog.Debug("successfully locked the hardware device", "camera", c, "driver", c.hardwareCamera)

	defer func() {
		c.hardwareCamera.mutex.Unlock()
		slog.Debug("successfully unlocked the hardware device", "camera", c, "driver", c.hardwareCamera)
	}()

	if !c.isRunning {
		return nil
	}

	c.isRunning = false
	c.hardwareCamera.cameraCounter--
	slog.Debug("decrement the hardware device counter", "camera", c, "driver", c.hardwareCamera, "counter", c.hardwareCamera.cameraCounter)

	c.hardwareCamera.stop()

	gstElementsSize := 2 + len(c.gstAdditionalElements)
	allGstElements := make([]*gst.Element, gstElementsSize)
	allGstElements[0] = c.gstQueueElement
	for i, element := range c.gstAdditionalElements {
		allGstElements[i+1] = element
	}
	allGstElements[gstElementsSize-1] = c.gstAppSink.Element

	ok := c.gstTeeSrcPad.Unlink(c.gstQueueSinkPad)
	if !ok {
		return fmt.Errorf("could not unlink tee and queue")
	}
	gst.ElementUnlinkMany(allGstElements...)

	err := c.hardwareCamera.gstPipeline.RemoveMany(allGstElements...)
	if err != nil {
		return err
	}

	c.hardwareCamera.gstTeeElement.ReleaseRequestPad(c.gstTeeSrcPad)

	if c.hardwareCamera.cameraCounter > 0 {
		c.hardwareCamera.play()
	} else {
		c.hardwareCamera.stop()
	}

	close(c.frameChannel)
	return nil
}

func (c *Camera) Frame() chan []byte {
	return c.frameChannel
}

func (c *Camera) onNewSample(sink *app.Sink) gst.FlowReturn {
	sample := sink.PullSample()
	if sample == nil {
		return gst.FlowEOS
	}

	// Retrieve the buffer from the sample
	buffer := sample.GetBuffer()
	if buffer == nil {
		return gst.FlowError
	}

	buffer.Map(gst.MapRead)
	frame := buffer.Bytes()
	defer buffer.Unmap()

	select {
	case c.frameChannel <- frame:
	default:
	}
	return gst.FlowOK
}
