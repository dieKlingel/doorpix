package oplog

import (
	"log/slog"

	"github.com/dieklingel/doorpix/internal/eventemitter"
)

type Event = eventemitter.Event
type Logger struct {
	emitter eventemitter.EventEmitter
	writer  Writer
}

var logger = &Logger{
	emitter: eventemitter.NewEventEmitter(),
}

func Default() *Logger {
	return logger
}

func (l *Logger) SetWriter(wr Writer) {
	l.writer = wr
}

func On(topic string) <-chan Event {
	return logger.emitter.On(topic)
}

func Dispatch(topic string, args ...any) {
	evt, err := logger.emitter.Dispatch(topic, args...)
	if err != nil {
		slog.Error("oplog: failed to dispatch event", "error", err.Error())
		return
	}

	if logger.writer != nil {
		err := logger.writer.Write(evt)
		if err != nil {
			slog.Error("oplog: failed to write event", "error", err.Error())
		}
	}

}

func Events() []Event {
	return logger.emitter.Events()
}
