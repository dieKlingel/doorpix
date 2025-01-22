package core

type Handler interface {
	Setup()
	Cleanup()
	Exec()
}
