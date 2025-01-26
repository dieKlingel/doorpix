package camera

import "github.com/go-gst/go-gst/gst"

type ElementFactory func() []*gst.Element

var FullHD ElementFactory = func() []*gst.Element {
	return []*gst.Element{
		MustNewElement("videoscale"),
		MustNewElement(
			"capsfilter",
			"caps", gst.NewCapsFromString("video/x-raw,width=1920,height=1080"),
		),
	}
}

var JPEG ElementFactory = func() []*gst.Element {
	return []*gst.Element{
		MustNewElement("jpegenc"),
	}
}

var I420 ElementFactory = func() []*gst.Element {
	return []*gst.Element{
		MustNewElement("videoconvert"),
		MustNewElement(
			"capsfilter",
			"caps", gst.NewCapsFromString("video/x-raw,format=I420,framerate=30/1"),
		),
	}
}

func New(identifier string, elementFactories ...ElementFactory) (*Camera, error) {
	allElements := make([]*gst.Element, 0)

	for _, createElements := range elementFactories {
		elements := createElements()
		allElements = append(allElements, elements...)
	}

	return NewFromString(identifier, allElements...)
}
