package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/dikaeinstein/go-microservice/chapter1/rpc/server"
)

type HTTPConn struct {
	in  io.Reader
	out io.Writer
}

func (c *HTTPConn) Read(p []byte) (int, error) {
	return c.in.Read(p)
}

func (c *HTTPConn) Write(p []byte) (int, error) {
	return c.out.Write(p)
}

func (c *HTTPConn) Close() error { return nil }

func Start() {
	port := 8080
	helloWorld := new(server.HelloWorldHandler)
	rpc.Register(helloWorld)

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Unable to listen to port: %v", err)
	}

	http.Serve(l, http.HandlerFunc(httpHander))
}

func httpHander(w http.ResponseWriter, r *http.Request) {
	serverCodec := jsonrpc.NewServerCodec(&HTTPConn{in: r.Body, out: w})
	err := rpc.ServeRequest(serverCodec)
	if err != nil {
		log.Printf("Error while serving JSON request: %v", err)
		http.Error(w, "Error while serving JSON request, details have been logged.",
			http.StatusInternalServerError)
		return
	}
}
