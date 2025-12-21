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
		http.Error(w, "webcam not configured", http.StatusInternalServerError)
		return
	}

	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *handler) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	slog.Debug("handle camera snapshot")

	if h.webcam == nil {
		http.Error(w, "webcam not configured", http.StatusInternalServerError)
		return
	}

	http.Error(w, "not implemented", http.StatusNotImplemented)
}
