package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
)

type HTTPService struct {
	Config doorpix.Config

	server *http.Server
}

type APIEventRequest struct {
	Event doorpix.EventType `json:"event"`
	Data  map[string]any    `json:"data"`
}

func (service *HTTPService) Init() error {
	slog.Debug("init http service")

	port := service.Config.HTTP.Port
	if port <= 0 {
		port = 8080
	}

	handler := http.NewServeMux()
	handler.HandleFunc("GET /api/camera/stream", service.showCameraStream)
	handler.HandleFunc("GET /api/camera/snapshot", service.showCameraFrame)

	// global config
	service.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	slog.Debug("successfully initialized http service")
	return nil
}

func (service *HTTPService) Exec(ctx context.Context, wg *sync.WaitGroup) error {
	slog.Debug("exec http service")

	go func() {
		if err := service.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		<-ctx.Done()
		slog.Debug("shutting down http service")

		if err := service.server.Shutdown(ctx); err != nil {
			slog.Warn("shut down http service with error", "error", err)
		} else {
			slog.Debug("successfully shut down http service")
		}
		wg.Done()
	}()

	return nil
}

func (h *HTTPService) showCameraStream(w http.ResponseWriter, r *http.Request) {
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

func (service *HTTPService) showCameraFrame(w http.ResponseWriter, r *http.Request) {
	webcam, err := service.newCamera()
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

func (service *HTTPService) newCamera() (*camera.Camera, error) {
	webcam, err := camera.New(
		service.Config.Camera.Device,
		camera.JPEG,
	)

	return webcam, err
}
