package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dikaeinstein/go-microservice/chapter4/data"
)

type searchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Kittens []data.Kitten `json:"kittens"`
}

type SearchHandler struct {
	DataStore data.Store
}

func (s SearchHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var request searchRequest
	err := decoder.Decode(&request)
	if err != nil || len(request.Query) < 1 {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	kittens, err := s.DataStore.Search(request.Query)
	if err != nil {
		http.Error(rw, fmt.Sprintf("Error fetching kittens: %v", err),
			http.StatusInternalServerError)
	}

	rw.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(rw)
	encoder.Encode(&SearchResponse{Kittens: kittens})
}
