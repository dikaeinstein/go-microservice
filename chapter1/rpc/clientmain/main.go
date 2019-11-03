package main

import (
	"fmt"

	rpcClient "github.com/dikaeinstein/go-microservice/chapter1/rpc/client"
	"github.com/dikaeinstein/go-microservice/chapter1/rpchttpjson/client"
)

func main() {
	rpcJSONClient := client.CreateClient()
	helloWorldResponse := rpcClient.PerformRequest(rpcJSONClient)

	fmt.Println(helloWorldResponse.Message)
}
