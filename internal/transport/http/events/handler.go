package events

import (
	"encoding/json"
	"net/http"

	"github.com/dieklingel/doorpix/internal/oplog"
	"github.com/gorilla/mux"
)

func Handler() http.Handler {
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

		args := make([]any, 0, len(body.Properties)*2+2)
		args = append(args, "source", "http")
		for key, value := range body.Properties {
			args = append(args, key, value)
		}

		oplog.Dispatch(body.Path, args...)
		w.WriteHeader(http.StatusCreated)
	})

	return router
}
