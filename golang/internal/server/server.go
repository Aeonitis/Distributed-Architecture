package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

// NewHTTPServer Takes in an address for the server to run on and returns an *http.Server.
// Create server with popular gorilla/mux library to write
// RESTful routes that match incoming requests to their respective handlers.
func NewHTTPServer(address string) *http.Server {
	httpServer := newHTTPServer()
	router := mux.NewRouter()

	// POST request to '/' matches the produce handler and appends the record to the log
	router.HandleFunc("/", httpServer.handleProduce).Methods("POST")
	// GET request to / matches the consume handler and reads the record from the log
	router.HandleFunc("/", httpServer.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    address,
		Handler: router,
	}
}
