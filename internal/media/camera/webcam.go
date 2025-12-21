package camera

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type Webcam struct {
	name   string
	driver Driver

	sessions map[string]*session
	mutex    sync.Mutex
}

func NewWebcam(name string, driver Driver) (*Webcam, error) {
	webcam := &Webcam{
		name:     name,
		driver:   driver,
		sessions: make(map[string]*session),
	}

	appsink, err := driver.GetAppSinkByName(name)
	if err != nil {
		return nil, err
	}

	appsink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: webcam.onNewSample,
	})

	return webcam, nil
}

func (w *Webcam) Start() (Session, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	slog.Debug("start webcam")

	sessionName := getName("session")
	if _, exists := w.sessions[sessionName]; exists {
		return nil, fmt.Errorf("duplicate key error")
	}

	err := w.driver.Start(w.name)
	if err != nil {
		return nil, err
	}

	session := &session{
		name:   sessionName,
		webcam: w,
		frame:  make(chan []byte),
	}
	w.sessions[sessionName] = session

	slog.Debug("start new session", "name", session.name, "webcam", w.name)
	return session, nil
}

func (w *Webcam) stop(session *session) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	slog.Debug("stop session", "session", session.name, "webcam", w.name)
	delete(w.sessions, session.name)

	if len(w.sessions) > 0 {
		slog.Debug("keep webcame active, as there are still active sessions")
		return nil
	}

	return w.driver.Stop(w.name)
}

func (w *Webcam) onNewSample(sink *app.Sink) gst.FlowReturn {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	sample := sink.PullSample()
	if sample == nil {
		return gst.FlowEOS
	}

	// Retrieve the buffer from the sample
	buffer := sample.GetBuffer()
	if buffer == nil {
		return gst.FlowError
	}

	buffer.Map(gst.MapRead)
	frame := buffer.Bytes()
	defer buffer.Unmap()

	for _, session := range w.sessions {
		select {
		case session.frame <- frame:
		default:
		}
	}

	return gst.FlowOK
}
