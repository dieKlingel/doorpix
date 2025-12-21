package main

import (
	"github.com/dieklingel/doorpix/internal/transport/http"
)

func main() {
	httpServer := http.NewServer(http.ServerProps{})
	httpServer.Serve()
}
