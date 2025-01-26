package core

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type HttpService struct {
	System doorpix.System

	server *http.Server
}

type APIEventRequest struct {
	Event doorpix.EventType `json:"event"`
	Data  map[string]any    `json:"data"`
}

func (service *HttpService) HandleEvent(config doorpix.Action, event *doorpix.Event) {
}

func (service *HttpService) Setup() {
	service.System.Bus.Handler(service)

	port := service.System.Config.HTTP.Port
	if port <= 0 {
		port = 8080
	}

	handler := http.NewServeMux()
	handler.HandleFunc("POST /api/events", service.AddNewEvent)
	handler.HandleFunc("GET /api/camera/stream", service.showCameraStream)
	handler.HandleFunc("GET /api/camera/snapshot", service.showCameraFrame)

	// global config
	service.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
}

func (service *HttpService) Exec() {
	go func() {
		if err := service.server.ListenAndServe(); err != nil {
			slog.Error("http server error", "error", err)
		}
	}()
}

func (service *HttpService) Cleanup() {}

func (service *HttpService) AddNewEvent(w http.ResponseWriter, r *http.Request) {
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
			service.System.Bus.OnWithData(req.Event, req.Data)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "event not allowed", http.StatusBadRequest)
}

func (h *HttpService) showCameraStream(w http.ResponseWriter, r *http.Request) {
	webcam, err := h.newCamera()
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

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	for {
		select {
		case <-r.Context().Done():
			return
		case frame, ok := <-webcam.Frame():
			if !ok {
				slog.Warn("no frame was received from the camera, canceling the stream")
				break
			}

			w.Write([]byte("--frame\n"))
			w.Write([]byte("Content-Type: image/jpeg\n\n"))
			w.Write(frame)
			w.Write([]byte("\n"))
		}
	}
}

func (h *HttpService) showCameraFrame(w http.ResponseWriter, r *http.Request) {
	webcam, err := h.newCamera()
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

func (service *HttpService) newCamera() (*camera.Camera, error) {
	webcam, err := camera.New(
		service.System.Config.Camera.Device,
		camera.JPEG,
	)

	return webcam, err
}
