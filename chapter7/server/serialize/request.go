package serialize

import (
	"encoding/json"
	"net/http"
	"strings"
)

// SerializableRequest wraps *http.Request
// and can also serialize to json
type SerializableRequest struct {
	*http.Request
}

// NewSerializableRequest creates a new instance of a SerializableRequest
func NewSerializableRequest(h *http.Request) *SerializableRequest {
	return &SerializableRequest{h}
}

// ToJSON serializes the *http.Request and returns as string
func (sr *SerializableRequest) ToJSON() string {
	data, _ := json.Marshal(sr.Serialize())
	return string(data)
}

func (sr *SerializableRequest) ToMap() map[string]interface{} {
	var headers []SerializedHeader
	for k, v := range sr.Header {
		headers = append(headers, SerializedHeader{
			Key: k, Value: strings.Join(v, ","),
		})
	}

	return map[string]interface{}{
		"method":     sr.Method,
		"host":       sr.Host,
		"query":      sr.URL.RawQuery,
		"path":       sr.URL.Path,
		"url":        sr.RemoteAddr,
		"headers":    headers,
		"statusCode": sr.Response.StatusCode,
		"error":      nil,
	}
}

// Serialize incoming *http.Request
func (sr *SerializableRequest) Serialize() SerializedRequest {
	var headers []SerializedHeader
	for k, v := range sr.Header {
		headers = append(headers, SerializedHeader{
			Key: k, Value: strings.Join(v, ","),
		})
	}

	return SerializedRequest{
		Method:  sr.Method,
		Host:    sr.Host,
		Query:   sr.URL.RawQuery,
		Path:    sr.URL.Path,
		Headers: headers,
	}
}

type SerializedRequest struct {
	Method  string
	Host    string
	Query   string
	Path    string
	Headers []SerializedHeader
}

type SerializedHeader struct {
	Key   string
	Value string
}
