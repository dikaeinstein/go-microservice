package main

import (
	"log"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8091"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
}
