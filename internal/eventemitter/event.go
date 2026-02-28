package eventemitter

import (
	"time"
)

type Event struct {
	Id         string         `json:"id"`
	Timestamp  time.Time      `json:"timestamp"`
	Path       string         `json:"path"`
	Properties map[string]any `json:"properties"`
}
