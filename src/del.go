package main

import (
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/msesto/cloudnative-go/src/handlers"
)

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
