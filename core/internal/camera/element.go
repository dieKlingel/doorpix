package camera

import (
	"github.com/go-gst/go-gst/gst"
)

func NewElement(name string, properties ...any) *gst.Element {
	props := make(map[string]any, 0)

	if len(properties)%2 != 0 {
		panic("invalid properties count")
	}

	for i := 0; i < len(properties); i += 2 {
		props[properties[i].(string)] = properties[i+1]
	}

	element, err := gst.NewElementWithProperties(name, props)
	if err != nil {
		panic(err)
	}

	return element
}
