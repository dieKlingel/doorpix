package doorpix

import (
	"os"
	"time"
)

type ActionHook struct {
	Data map[string]any
}

func NewActionHook(data map[string]any) *ActionHook {
	event := &ActionHook{
		Data: data,
	}

	if event.Data == nil {
		event.Data = make(map[string]any)
	}

	event.Data["date"] = time.Now().Format("2006-01-02")
	event.Data["time"] = time.Now().Format("15:04:05")
	if dir, err := os.Getwd(); err == nil {
		event.Data["pwd"] = dir
	}

	return event
}
