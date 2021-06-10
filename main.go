package main

import (
	"errors"
	"fmt"
	"io"
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

func KeyValuePutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := io.ReadAll(r.Body) // The request body has our value
	defer r.Body.Close()

	if err != nil { // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}

	err = h.Put(key, string(value)) // Store the value as a string
	if err != nil {                 // If we have an error, report it
		http.Error(w,
			err.Error(),
			http.StatusInternalServerError)
		return
	}
	logger.WritePut(key, string(value))
	w.WriteHeader(http.StatusCreated) // All good! Return StatusCreated
}

func KeyValueDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	err := h.Delete(key) // Delete key

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.WriteDelete(key)
	w.WriteHeader(http.StatusOK)
}

func KeyValueGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Retrieve "key" from the request
	key := vars["key"]

	value, err := h.Get(key) // Get value for key
	if errors.Is(err, h.ErrorNoSuchKey) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(value)) // Write the value to the response
}
