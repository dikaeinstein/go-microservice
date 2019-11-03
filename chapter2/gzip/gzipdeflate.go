package main

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	gzw *gzip.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	contentType := grw.Header().Get("Content-Type")
	if contentType == "" {
		// If content type is not set, infer it from the uncompressed body.
		grw.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return grw.gzw.Write(b)
}

type deflateResponseWriter struct {
	dw *flate.Writer
	http.ResponseWriter
}

func (drw deflateResponseWriter) Write(b []byte) (int, error) {
	contentType := drw.Header().Get("Content-Type")
	if contentType == "" {
		// If content type is not set, infer it from the uncompressed body.
		drw.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return drw.dw.Write(b)
}

type gzipHandler struct {
	next http.Handler
}

func (h *gzipHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	contentEncoding := r.Header.Get("Accept-Encoding")
	rw.Header().Set("Content-Type", "application/json")
	if strings.Contains(contentEncoding, "gzip") {
		h.ServeGzipped(rw, r)
	} else if strings.Contains(contentEncoding, "deflate") {
		h.ServeDeflated(rw, r)
	} else {
		h.ServePlain(rw, r)
	}
}

func (h *gzipHandler) ServeGzipped(rw http.ResponseWriter, r *http.Request) {
	gzw := gzip.NewWriter(rw)
	defer gzw.Close()

	rw.Header().Set("Content-Encoding", "gzip")
	h.next.ServeHTTP(gzipResponseWriter{gzw, rw}, r)
}

func (h *gzipHandler) ServeDeflated(rw http.ResponseWriter, r *http.Request) {
	dw, err := flate.NewWriter(rw, flate.DefaultCompression)
	if err != nil {
		log.Fatal(err)
	}
	defer dw.Close()

	rw.Header().Set("Content-Encoding", "deflate")
	h.next.ServeHTTP(deflateResponseWriter{dw, rw}, r)
}

func (h *gzipHandler) ServePlain(w http.ResponseWriter, r *http.Request) {
	h.next.ServeHTTP(w, r)
}

type helloWorldRequest struct {
	Name string `json:"name"`
}

type helloWorldResponse struct {
	Message string `json:"message"`
}

type validatedContextKey string

func main() {
	port := 8080
	gzipHelloWorldHandler := &gzipHandler{http.HandlerFunc(helloWorldHandler)}
	http.Handle("/helloworld", gzipHelloWorldHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func helloWorldHandler(wr http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(wr, "Bad request", http.StatusBadRequest)
		return
	}

	response := helloWorldResponse{Message: "Hello " + request.Name}

	encoder := json.NewEncoder(wr)
	encoder.Encode(response)
}
