package handler

import "net/http"

type bang struct{}

// NewBangHandler creates a handler which panics
func NewBangHandler() http.Handler {
	return &bang{}
}

func (bang) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("Somethings gone wrong again")
}
