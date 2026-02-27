package events

import (
	"encoding/json"
	"net/http"

	"github.com/dieklingel/doorpix/internal/eventemitter"
	"github.com/gorilla/mux"
)

func Handler(oplog eventemitter.EventEmitter) http.Handler {
	router := mux.NewRouter()

	router.Path("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		evts := oplog.Events()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(evts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	router.Path("/").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body CreateEventRequest
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		event, err := oplog.DispatchProperties(body.Path, body.Properties)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(event)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return router
}
