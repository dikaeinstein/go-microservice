package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dikaeinstein/go-microservice/chapter4/data"
	"github.com/dikaeinstein/go-microservice/chapter4/handlers"
)

func main() {
	store, err := data.NewMongoStore("mongodb://localhost:27017/kittenserver")
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":2323", &handlers.SearchHandler{DataStore: store})
	if err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("Server closed.")
			os.Exit(0)
		} else {
			log.Fatal(err)
		}
	}
}
