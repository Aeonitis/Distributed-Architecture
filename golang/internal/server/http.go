package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/** Handlers will:
1. Unmarshal the request’s JSON body into a struct.
2. Run that endpoint’s logic with the request to obtain a result.
3. Marshal and write that result to the response.
*/

// httpServer Server referencing a log to defer to in its handlers
type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

// ProduceRequest contains the record that the caller of our API wants appended to the log
type ProduceRequest struct {
	Record Record `json:"record"`
}

// ProduceResponse contains the record that the caller of our API wants appended to the log
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

// ConsumeRequest specifies which record caller of our API wants to read
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

// ConsumeResponse to send back records to the caller of our API
type ConsumeResponse struct {
	Record Record `json:"record"`
}

// handleProduce implements:
// 	I. Unmarshal request’s JSON body into a struct
// 	II. Produce log and get offset that record is stored under
// 	III. Marshal and write that result to the response
func (server *httpServer) handleProduce(w http.ResponseWriter, httpRequest *http.Request) {
	var request ProduceRequest
	fmt.Println("HandleProduce: ", httpRequest.Body)

	err := json.NewDecoder(httpRequest.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := server.Log.Append(request.Record)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleConsume implements:
// 	I. Unmarshal request’s JSON body into a struct
// 	II. Produce log by calling Read(offset uint64) to get record stored in the log.
// 	III. Marshal and write that result to the response
func (server *httpServer) handleConsume(writer http.ResponseWriter, httpRequest *http.Request) {
	var request ConsumeRequest
	fmt.Println("HandleConsume: ", httpRequest.Body)

	err := json.NewDecoder(httpRequest.Body).Decode(&request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := server.Log.Read(request.Offset)

	fmt.Println("Current Record: ", record)

	// More error checking to provide an accurate status code to the client if server can’t handle the request,
	// e.g. if the client requested a record that does not exist.
	if err == ErrMessageOffsetNotFound {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	response := ConsumeResponse{Record: record}
	err = json.NewEncoder(writer).Encode(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
