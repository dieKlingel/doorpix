package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type HttpHandler struct {
	System doorpix.System

	server *http.Server
}

type APIEventRequest struct {
	Event doorpix.EventType `json:"event"`
	Data  map[string]any    `json:"data"`
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
	handler.HandleFunc("POST /api/events", h.AddNewEvent)
	handler.HandleFunc("GET /api/camera/stream", h.showCameraStream)
	handler.HandleFunc("GET /api/camera/snapshot", h.showCameraFrame)

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

func (h *HttpHandler) AddNewEvent(w http.ResponseWriter, r *http.Request) {
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
			h.System.Bus.OnWithData(req.Event, req.Data)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "event not allowed", http.StatusBadRequest)
}

func (h *HttpHandler) showCameraStream(w http.ResponseWriter, r *http.Request) {
	webcam, err := camera.NewFromString(
		h.System.Config.Camera.Device,
		camera.MustNewElement("jpegenc"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Debug("b starting webcam")
	err = webcam.Start()
	slog.Debug("b webcam started")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		slog.Info("stopping webcam")
		err = webcam.Stop()
		if err != nil {
			slog.Error("error stopping webcam", "error", err)
		}
	}()

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	for {
		select {
		case <-r.Context().Done():
			slog.Info("stream closed")
			return
		case frame, ok := <-webcam.Frame():
			if !ok {
				slog.Warn("no frame")
				break
			}

			w.Write([]byte("--frame\n"))
			w.Write([]byte("Content-Type: image/jpeg\n\n"))
			w.Write(frame)
			w.Write([]byte("\n"))
		}
	}
}

func (h *HttpHandler) showCameraFrame(w http.ResponseWriter, r *http.Request) {
	webcam, err := camera.NewFromString(
		h.System.Config.Camera.Device,
		camera.MustNewElement("jpegenc"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = webcam.Start()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		err = webcam.Stop()
		if err != nil {
			slog.Error("error stopping webcam", "error", err)
		}
	}()

	select {
	case <-r.Context().Done():
		return
	case frame, ok := <-webcam.Frame():
		if !ok {
			http.Error(w, "no frame", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(frame)
	}
}
