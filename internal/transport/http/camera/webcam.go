package camera

type Webcam interface {
	Start() error
	Stop() error
	Frame() chan []byte
}
