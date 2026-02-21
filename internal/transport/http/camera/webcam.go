package camera

import "github.com/dieklingel/doorpix/internal/media/camera"

type Webcam interface {
	Start() (camera.Session, error)
}
