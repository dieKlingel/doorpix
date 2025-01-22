package doorpix

import (
	"os"
	"time"
)

type Event struct {
	Data map[string]interface{}
}

func NewEvent() *Event {
	event := &Event{
		Data: make(map[string]interface{}),
	}

	event.Data["date"] = time.Now().Format("2006-01-02")
	event.Data["time"] = time.Now().Format("15:04:05")
	if dir, err := os.Getwd(); err == nil {
		event.Data["pwd"] = dir
	}

	return event
}
