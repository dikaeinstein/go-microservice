package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/dikaeinstein/go-microservice/chapter1/rpc/server"
)

func Start() {
	port := 8080
	helloWorld := &server.HelloWorldHandler{}
	rpc.Register(helloWorld)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Unable to listen on port: %v\n", err)
	}

	log.Printf("Server starting on port: %v\n", port)

	http.Serve(l, nil)
}
