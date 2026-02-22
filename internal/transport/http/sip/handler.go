package sip

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/dieklingel/doorpix/internal/transport/sip"
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
	router.Path("/calls").Methods("POST").HandlerFunc(handler.createCall)
	router.Path("/calls/{id}").Methods("GET").HandlerFunc(handler.getCallById)
	router.Path("/calls/{id}").Methods("DELETE").HandlerFunc(handler.hangupCall)
	router.Path("/account").Methods("GET").HandlerFunc(handler.getAccount)

	return router
}

func (h *handler) getAllCalls(w http.ResponseWriter, r *http.Request) {
	calls := h.ua.Calls()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(calls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *handler) getCallById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	call := h.ua.CallById(id)

	w.Header().Set("Content-Type", "application/json")
	if call == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(call)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) hangupCall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.ua.Hangup(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusGone)
}

func (h *handler) getAccount(w http.ResponseWriter, r *http.Request) {
	info := h.ua.AccountInfo()
	var payload any

	if info != nil {
		payload = info
	} else {
		payload = struct{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) createCall(w http.ResponseWriter, r *http.Request) {
	var body CreateCallRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	call, err := h.ua.Invite(body.Uri)
	if err != nil {
		if errors.Is(err, sip.ErrInvalidUri) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else if errors.Is(err, sip.ErrNotReady) {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(call)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
