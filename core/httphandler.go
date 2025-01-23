package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type HttpHandler struct {
	System doorpix.System

	server *http.Server
}

type APIEventRequest struct {
	Event doorpix.EventType `json:"event"`
}

func (h *HttpHandler) HandleEvent(config doorpix.Action, event *doorpix.Event) {
}

func (h *HttpHandler) Setup() {
	h.System.Bus.Handler(h)

	port := h.System.Config.HTTP.Port
	if port <= 0 {
		port = 8080
	}

	handler := http.NewServeMux()
	handler.HandleFunc("POST /api/events", h.handleEmitEvent)

	// global config
	h.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}

func (h *HttpHandler) Exec() {
	go func() {
		if err := h.server.ListenAndServe(); err != nil {
			slog.Error("http server error", "error", err)
		}
	}()
}

func (h *HttpHandler) Cleanup() {}

func (h *HttpHandler) handleEmitEvent(w http.ResponseWriter, r *http.Request) {
	var req APIEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allowedEventTypes := []doorpix.EventType{
		doorpix.APIRingEvent,
		doorpix.APIUnlockEvent,
	}

	for _, eventType := range allowedEventTypes {
		if req.Event == eventType {
			h.System.Bus.On(req.Event)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "event not allowed", http.StatusBadRequest)
}
