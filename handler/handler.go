package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhampac/shortnr/storage"
)

// type response struct {
// 	Success bool        `json:"success"`
// 	Data    interface{} `json:"short_url"`
// }

// type handler struct {
// 	schema  string
// 	host    string
// 	storage storage.Service
// }

// func (h handler) encode(ctx context.Context) (interface{}, int, error) {
// 	var input struct {
// 		URL     string `json:"url"`
// 		Expires string `json:"expires"`
// 	}

// 	return nil, 0, nil
// }

func New(schema string, host string, storage storage.Service) *mux.Router {
	r := mux.NewRouter()
	// h := handler{schema, host, storage}

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/encode", encodeHandler).Methods("POST")
	return r
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the best URL shortener service!")
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Coming soon")
}
