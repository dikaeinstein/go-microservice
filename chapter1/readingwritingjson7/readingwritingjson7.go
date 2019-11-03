package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type helloWorldRequest struct {
	Name string `json:"name"`
}

type helloWorldResponse struct {
	Message string `json:"message"`
}

type validationHandler struct {
	next http.Handler
}

func newValidationHandler(next http.Handler) http.Handler {
	return validationHandler{next: next}
}

func (h validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
		return
	}

	h.next.ServeHTTP(rw, r)
}

func main() {
	port := 8080

	http.HandleFunc("/helloworld", helloWorldHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(wr http.ResponseWriter, r *http.Request) {
	name := r.Context().Value("name").(string)
	response := helloWorldResponse{Message: "Hello " + name}

	encoder := json.NewEncoder(wr)
	encoder.Encode(response)
}
