package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dikaeinstein/go-microservice/chapter4/data"
)

func BenchmarkSearchHandler(b *testing.B) {
	mockStore := &data.MockStore{}
	mockStore.On("Search", "Fat Freddy's Cat").Return([]data.Kitten{
		{
			ID:     "2",
			Name:   "Fat Freddy's Cat",
			Weight: 20.0,
		},
	})

	s := SearchHandler{DataStore: mockStore}

	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(http.MethodGet, "/search", bytes.NewReader([]byte(`{"query":"Fat Freddy's Cat"}`)))
		rr := httptest.NewRecorder()
		s.ServeHTTP(rr, r)
	}
}
