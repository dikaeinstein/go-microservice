package server

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/dikaeinstein/go-microservice/chapter1/rpc/contract"
)

type HelloWorldHandler struct{}

func (h *HelloWorldHandler) HelloWorld(
	args *contract.HelloWorldRequest, reply *contract.HelloWorldResponse,
) error {
	reply.Message = "hello " + args.Name
	return nil
}

func StartServer() {
	port := 8080
	helloWorld := &HelloWorldHandler{}

	rpc.Register(helloWorld)
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("Unable to listen on given port: %s", err)
	}

	for {
		conn, _ := l.Accept()
		go rpc.ServeConn(conn)
	}
}
