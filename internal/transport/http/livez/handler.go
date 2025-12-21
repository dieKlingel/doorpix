package livez

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Handler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK\r\n"))
	})

	return router
}
