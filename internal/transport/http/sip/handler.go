package sip

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type handler struct {
	ua UserAgent
}

func Handler(ua UserAgent) http.Handler {
	router := mux.NewRouter()

	handler := &handler{
		ua: ua,
	}

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\r\n"))
	})
	router.Path("/calls").Methods("GET").HandlerFunc(handler.getAllCalls)
	router.Path("/account").Methods("GET").HandlerFunc(handler.getAccount)

	return router
}

func (h *handler) getAllCalls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("[]\r\n"))
}

func (h *handler) getAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	info := h.ua.AccountInfo()
	var payload any

	if info != nil {
		payload = info
	} else {
		payload = struct{}{}
	}

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
