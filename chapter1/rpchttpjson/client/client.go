package client

import (
	"fmt"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func CreateClient() *rpc.Client {
	port := 8080
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal("Error Dialing: ", err)
	}

	return client
}
