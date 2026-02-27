package eventemitter

import "time"

type Event struct {
	Timestamp  time.Time      `json:"timestamp"`
	Path       string         `json:"path"`
	Properties map[string]any `json:"properties"`
}
