package camera

import (
	"fmt"
	"math/rand/v2"
)

type Session struct {
	name   string
	webcam *Webcam
	frame  chan []byte
}

func (session *Session) Stop() error {
	session.webcam.stop(session)
	close(session.frame)

	return nil
}

func (session *Session) Frame() chan []byte {
	return session.frame
}

func getName(n string) string {
	suffix := rand.IntN(8999)
	suffix = suffix + 1000 // 1000 <= suffx <= 9999
	return fmt.Sprintf("%s-%d", n, suffix)
}
