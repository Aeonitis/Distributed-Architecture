package main

import (
	"example.com/m/v2/internal/server"
	"log"
)

// main creates and start the server, passing in the address to listen on (localhost:8080)
// Wrapping our server with the *net/http.Server in NewHTTPServer() for a quick web server
func main() {
	simpleHttpServer := server.NewHTTPServer(":8080")

	// Invoke server to listen for and handle requests by calling ListenAndServe()
	log.Fatal(simpleHttpServer.ListenAndServe())
}

/**
------------ TEST----------->

Produce Resources to append/set values of new records
$ curl -X POST localhost:8080 -d '{"record": {"value": "Emil"}}'
$ curl -X POST localhost:8080 -d '{"record": {"value": "1.22474487139..."}}'
$ curl -X POST localhost:8080 -d '{"record": {"value": "Kaine"}}'

Consume Resources to read/check records by index
$ curl -X GET localhost:8080 -d '{"offset": 0}'
$ curl -X GET localhost:8080 -d '{"offset": 1}'
$ curl -X GET localhost:8080 -d '{"offset": 2}'
*/
