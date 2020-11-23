package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
)

type ExpandRequest struct {
	Query             string               `json:"query"`
	AddressComponents []string             `json:"address_components,omitempty"`
	Options           expand.ExpandOptions `json:"options,omitempty"`
}

type ParseRequest struct {
	Query   string               `json:"query"`
	Options parser.ParserOptions `json:"options,omitempty"`
}

func GetAddressComponent(componentName string) (uint16, error) {
	var err error

	var result uint16 = expand.AddressNone
	switch strings.ToLower(componentName) {
	case "none":
		result = expand.AddressNone
	case "any":
		result = expand.AddressAny
	case "name":
		result = expand.AddressName
	case "house_number":
		result = expand.AddressHouseNumber
	case "street":
		result = expand.AddressStreet
	case "unit":
		result = expand.AddressUnit
	case "level":
		result = expand.AddressLevel
	case "staircase":
		result = expand.AddressStaircase
	case "entrance":
		result = expand.AddressEntrance
	case "category":
		result = expand.AddressCategory
	case "near":
		result = expand.AddressNear
	case "toponym":
		result = expand.AddressToponym
	case "postal_code":
		result = expand.AddressPostalCode
	case "po_box":
		result = expand.AddressPoBox
	case "all":
		result = expand.AddressAll
	default:
		err = fmt.Errorf("Unknown componentName %s", componentName)
	}

	return result, err
}

func ConstructAddressComponentsValue(component_names []string) uint16 {
	var result uint16 = 0
	for _, component_name := range component_names {
		component, _ := GetAddressComponent(component_name)
		result |= component
	}

	return result
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

	req := ExpandRequest{
		Options: expand.GetDefaultExpansionOptions(),
	}

	q, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(q, &req)

	if len(req.AddressComponents) > 0 {
		req.Options.AddressComponents = ConstructAddressComponentsValue(req.AddressComponents)
	}

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
