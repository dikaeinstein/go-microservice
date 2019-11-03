package main

import (
	"log"
	"net/http"

	"github.com/dikaeinstein/go-microservice/chapter9/CQRS/data"
	nc "github.com/dikaeinstein/go-microservice/chapter9/CQRS/nats"
	"github.com/dikaeinstein/go-microservice/chapter9/CQRS/product/read"
	"github.com/dikaeinstein/go-microservice/chapter9/CQRS/product/write"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func main() {
	readStore := data.SetupDB(read.Schema)
	writeStore := data.SetupDB(write.Schema)
	write.SeedDB(writeStore)

	readNC, err := nc.ConnectNATS(nats.DefaultURL, "Products Reader Subscriber")
	writeNC, err := nc.ConnectNATS(nats.DefaultURL, "Products Writer Publisher")

	if err != nil {
		log.Fatal(err)
	}
	defer readNC.Close()
	defer writeNC.Close()

	readNC.Subscribe("product.inserted", read.MakeProductMessageCallBack(readStore))

	writeHandler := write.NewHandler(writeStore, writeNC)
	readHandler := read.NewHandler(readStore)
	r := mux.NewRouter()
	r.HandleFunc("/products", readHandler.GetProducts).Methods(http.MethodGet)
	r.HandleFunc("/products", writeHandler.InsertProduct).Methods(http.MethodPost)
	r.HandleFunc("/stock", writeHandler.StockCount)

	log.Println("Starting products service on port 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
