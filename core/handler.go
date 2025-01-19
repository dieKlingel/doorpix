package core

type Handler interface {
	Setup(*EventEmitter)
	Cleanup()
	Exec()
}
