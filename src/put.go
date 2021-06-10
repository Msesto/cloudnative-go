package main

import (
	"io"
	"net/http"

	h "github.com/Msesto/cloudnative-go/src/handlers"
	"github.com/gorilla/mux"
)

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
