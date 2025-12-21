package camera

import (
	"github.com/go-gst/go-gst/gst/app"
)

type Driver interface {
	Start(name string) error
	Stop(name string) error
	GetAppSinkByName(name string) (*app.Sink, error)
}
