package client

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/dikaeinstein/go-microservice/chapter1/rpc/contract"
)

func CreateClient() *rpc.Client {
	port := 8080
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		log.Fatal("dialling: ", err)
	}

	return client
}

func PerformRequest(client *rpc.Client) contract.HelloWorldResponse {
	args := contract.HelloWorldRequest{Name: "Dika"}
	var reply contract.HelloWorldResponse
	err := client.Call("HelloWorldHandler.HelloWorld", args, &reply)
	if err != nil {
		log.Fatal("error client: ", err)
	}
	client.Close()
	return reply
}
