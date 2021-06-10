package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/msesto/cloudnative-go/src/handlers"
)

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
