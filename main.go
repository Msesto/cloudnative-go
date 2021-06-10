package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/msesto/cloudnative-go/handlers"
)

func helloMuxHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello gorilla/mux!\n"))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", helloMuxHandler)

	r.HandleFunc("/v1/{key}", h.KeyValuePutHandler).Methods("PUT")
	r.HandleFunc("/v1/{key}", h.KeyValueGetHandler).Methods("GET")
	r.HandleFunc("/v1/{key}", h.KeyValueDeleteHandler).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
