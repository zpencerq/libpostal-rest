package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
)

type ExpandRequest struct {
	Query string `json:"query"`
	Options expand.ExpandOptions `json:"options,omitempty"`
}

type ParseRequest struct {
	Query string `json:"query"`
	Options parser.ParserOptions `json:"options,omitempty"`
}

func main() {
	host := os.Getenv("LISTEN_HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "8080"
	}
	listenSpec := fmt.Sprintf("%s:%s", host, port)

	router := mux.NewRouter()
	router.HandleFunc("/health", HealthHandler).Methods("GET")
	router.HandleFunc("/expand", ExpandHandler).Methods("POST")
	router.HandleFunc("/parser", ParserHandler).Methods("POST")

	s := &http.Server{Addr: listenSpec, Handler: router}
	fmt.Printf("Listening on http://%s\n", listenSpec)
	s.ListenAndServe()
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func ExpandHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ExpandRequest

	q, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(q, &req)

	expansions := expand.ExpandAddressOptions(req.Query, req.Options)

	expansionThing, _ := json.Marshal(expansions)
	w.Write(expansionThing)
}

func ParserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ParseRequest

	q, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(q, &req)

	parsed := parser.ParseAddressOptions(req.Query, req.Options)
	parseThing, _ := json.Marshal(parsed)
	w.Write(parseThing)
}
