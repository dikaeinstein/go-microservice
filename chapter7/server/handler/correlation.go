package handler

import (
	"net/http"

	"github.com/google/uuid"
)

type tagRequest struct {
	next http.Handler
}

// NewTagRequest is a middleware handler which appends X-Request-ID
// header if one is not already set
func NewTagRequest(next http.Handler) http.Handler {
	return tagRequest{next}
}

func (t tagRequest) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Request-ID") == "" {
		r.Header.Set("X-Request-ID", uuid.New().String())
	}
	t.next.ServeHTTP(w, r)
}
