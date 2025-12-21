package camera

import (
	"fmt"
	"math/rand/v2"
)

type Session interface {
	Stop() error
	Frame() chan []byte
}

type session struct {
	name   string
	webcam *Webcam
	frame  chan []byte
}

func (session *session) Stop() error {
	session.webcam.stop(session)
	close(session.frame)

	return nil
}

func (session *session) Frame() chan []byte {
	return session.frame
}

func getName(n string) string {
	suffix := rand.IntN(8999)
	suffix = suffix + 1000 // 1000 <= suffx <= 9999
	return fmt.Sprintf("%s-%d", n, suffix)
}
