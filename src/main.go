package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/msesto/cloudnative-go/handlers"
)

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

var logger h.TransactionLogger

func initializeTransactionLog() error {
	var err error

	logger, err = h.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := h.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors: // Retrieve any errors
		case e, ok = <-events:
			switch e.EventType {
			case h.EventDelete: // Got a DELETE event!
				err = h.Delete(e.Key)
			case h.EventPut: // Got a PUT event!
				err = h.Put(e.Key, e.Value)
			}
		}
	}
	logger.Run()

	return err
}

func main() {
	initializeTransactionLog()
	r := mux.NewRouter()

	r.HandleFunc("/", helloMuxHandler)

	r.HandleFunc("/v1/{key}", KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", KeyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
