package main

import (
	"context"
	"log"

	kittens "github.com/dikaeinstein/go-microservice/chapter6/gomicro/proto"
	"github.com/micro/go-micro/config/cmd"
	"github.com/micro/go-micro/server"
)

func main() {
	cmd.Init()

	server.Init(
		server.Name("bmigo.micro.Kittens"),
		server.Version("1.0.0"),
		server.Address(":8091"),
	)

	// Register Handlers
	server.Handle(
		server.NewHandler(new(Kittens)),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

type Kittens struct{}

func (k Kittens) Hello(ctx context.Context, req *kittens.Request, rsp *kittens.Response) error {
	rsp.Msg = server.DefaultId + ": Hello " + req.Name

	return nil
}
