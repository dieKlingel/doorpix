package camera

import (
	"fmt"

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

func New(hardwareCamera *HardwareCamera, elements ...*gst.Element) (*Camera, error) {
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

func (c *Camera) Start() error {
	c.hardwareCamera.mutex.Lock()
	defer c.hardwareCamera.mutex.Unlock()
	if c.isRunning {
		return nil
	}

	c.hardwareCamera.cameraCounter++
	c.isRunning = true

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
	c.hardwareCamera.mutex.Lock()
	defer c.hardwareCamera.mutex.Unlock()
	if !c.isRunning {
		return nil
	}

	c.hardwareCamera.cameraCounter--
	c.isRunning = false

	c.hardwareCamera.pause()

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

	select {
	case c.frameChannel <- frame:
	default:
	}

	buffer.Unmap()

	return gst.FlowOK
}
