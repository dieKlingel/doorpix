package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dieklingel/doorpix/core/internal/camera"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/dieklingel/doorpix/core/internal/service"
)

type HTTPServerProps struct {
	Port                    int
	VideoStreamCameraDevice string
}

type HTTPServer struct {
	props HTTPServerProps

	server *http.Server
	ctx    service.Context
}

func NewHTTPServer(props HTTPServerProps) *HTTPServer {
	return &HTTPServer{
		props: props,
	}
}

type APIEventRequest struct {
	Event doorpix.EventType `json:"event"`
	Data  map[string]any    `json:"data"`
}

func (server *HTTPServer) Start() {
	server.ctx = service.NewContext(context.Background())

	port := server.props.Port
	if port <= 0 {
		port = 8080
	}

	handler := http.NewServeMux()
	handler.HandleFunc("GET /api/camera/stream", server.showCameraStream)
	handler.HandleFunc("GET /api/camera/snapshot", server.showCameraFrame)

	// global config
	server.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	server.StartBackgroundTask()
}

func (server *HTTPServer) Stop() {
	slog.Debug("stopping http service")

	server.ctx.CancelAndWait()
}

func (service *HTTPServer) StartBackgroundTask() {
	slog.Debug("run http service in background")

	go func() {
		if err := service.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("http server error", "error", err)
		}
	}()

	service.ctx.Lock()
	go func() {
		defer service.ctx.Unlock()

		<-service.ctx.Done()
		slog.Debug("shutting down http service")

		if err := service.server.Shutdown(service.ctx.Context()); err != nil {
			slog.Warn("shut down http service with error", "error", err)
		} else {
			slog.Debug("successfully shut down http service")
		}
	}()
}

func (h *HTTPServer) showCameraStream(w http.ResponseWriter, r *http.Request) {
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

func (service *HTTPServer) showCameraFrame(w http.ResponseWriter, r *http.Request) {
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

func (service *HTTPServer) newCamera() (*camera.Camera, error) {
	webcam, err := camera.New(
		service.props.VideoStreamCameraDevice,
		camera.JPEG,
	)

	return webcam, err
}
