package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dikaeinstein/go-microservice/chapter4/data"
)

var mockStore *data.MockStore

func setupTest(d interface{}) (*http.Request, *httptest.ResponseRecorder, SearchHandler) {
	mockStore = &data.MockStore{}
	s := SearchHandler{DataStore: mockStore}
	rw := httptest.NewRecorder()

	if d == nil {
		return httptest.NewRequest(http.MethodGet, "/search", nil), rw, s
	}

	body, _ := json.Marshal(d)
	return httptest.NewRequest(http.MethodGet, "/search", bytes.NewReader(body)), rw, s
}

func TestSearchHandlerReturnsBadRequestWhenNoSearchCriteriaIsSent(t *testing.T) {
	request, response, handler := setupTest(&searchRequest{})
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v", response.Code)
	}
}

func TestSearchHandlerReturnsBadRequestWhenBlankSearchCriteriaIsSent(t *testing.T) {
	r, rw, handler := setupTest(&searchRequest{Query: ""})
	handler.ServeHTTP(rw, r)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v", rw.Code)
	}
}

func TestSearchHandlerCallsDataStoreWithValidQuery(t *testing.T) {
	r, rw, handler := setupTest(&searchRequest{Query: "Fat Freddy's Cat"})
	mockStore.On("Search", "Fat Freddy's Cat").Return(make([]data.Kitten, 0))

	handler.ServeHTTP(rw, r)

	var response SearchResponse
	json.Unmarshal(rw.Body.Bytes(), &response)

	mockStore.AssertExpectations(t)
	assert.Len(t, response.Kittens, 0)
	assert.Equal(t, rw.Code, http.StatusOK)
}
