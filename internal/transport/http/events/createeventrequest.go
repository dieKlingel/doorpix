package events

type CreateEventRequest struct {
	Path       string         `json:"path"`
	Properties map[string]any `json:"properties"`
}
