package camera

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
	webcam Webcam
}

func Handler(webcam Webcam) http.Handler {
	router := mux.NewRouter()

	handler := &handler{webcam: webcam}

	router.HandleFunc("/stream", handler.handleStream).Methods("GET")
	router.HandleFunc("/snapshot", handler.handleSnapshot).Methods("GET")

	return router
}

func (h *handler) handleStream(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handle camera stream")

	if h.webcam == nil {
		http.Error(w, "camera not configured", http.StatusInternalServerError)
		return
	}

	session, err := h.webcam.Start()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		err = session.Stop()
		if err != nil {
			slog.Error("error stopping webcam", "error", err)
		}
	}()

	w.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")
	for {
		select {
		case <-r.Context().Done():
			return
		case frame, ok := <-session.Frame():
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

func (h *handler) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handle camera snapshot")

	if h.webcam == nil {
		http.Error(w, "camera not configured", http.StatusInternalServerError)
		return
	}

	session, err := h.webcam.Start()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		err = session.Stop()
		if err != nil {
			slog.Error("error stopping webcam", "error", err)
		}
	}()

	select {
	case <-r.Context().Done():
		return
	case frame, ok := <-session.Frame():
		if !ok {
			http.Error(w, "no frame", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(frame)
	}
}
