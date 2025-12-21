package http

import (
	"fmt"
	"net/http"

	"github.com/dieklingel/doorpix/internal/transport/http/camera"
	"github.com/dieklingel/doorpix/internal/transport/http/healthz"
	"github.com/dieklingel/doorpix/internal/transport/http/livez"
	"github.com/gorilla/mux"
)

const DefaultPort = 8080

type ServerProps struct {
	Port   *int
	Webcam *camera.Webcam
}

type Server struct {
	router *mux.Router
	port   int
}

func NewServer(props ServerProps) Server {
	router := mux.NewRouter()

	port := DefaultPort
	if props.Port != nil {
		port = *props.Port
	}

	server := Server{router: router, port: port}

	server.Subrouter("/livez", livez.Handler())
	server.Subrouter("/healthz", healthz.Handler())
	server.Subrouter("/camera", camera.Handler(props.Webcam))

	return server
}

func (s *Server) Serve() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	return server.ListenAndServe()
}

/**
 * Subrouter adds a subrouter to the server's main router at the specified path,
 * using the provided handler. It strips the prefix from incoming requests before
 * passing them to the handler.
 */
func (s *Server) Subrouter(path string, handler http.Handler) {
	subrouter := s.router.PathPrefix(path).Subrouter()
	subrouter.NewRoute().Handler(StripPrefix(path, handler))
}

/**
 * StripPrefix is a helper function that strips the specified prefix from the request URL
 * before passing the request to the provided handler. It ensures that the URL path
 * and raw path are set to "/" if they become empty after stripping the prefix.
 */
func StripPrefix(prefix string, h http.Handler) http.Handler {
	return http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}

		if r.URL.RawPath == "" {
			r.URL.RawPath = "/"
		}

		h.ServeHTTP(w, r)
	}))
}
