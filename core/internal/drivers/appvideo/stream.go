package appvideo

type Stream interface {
	Start() error
	Stop() error
	Frame() chan []byte
}
