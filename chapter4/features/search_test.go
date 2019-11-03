package features

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/dikaeinstein/go-microservice/chapter4/data"
	"github.com/dikaeinstein/go-microservice/chapter4/handlers"
)

var criteria interface{}
var err error
var response *http.Response

func iHaveNoSearchCriteria() error {
	if criteria != nil {
		return fmt.Errorf("criteria should be nil")
	}
	return nil
}

func iCallTheSearchEndpoint() error {
	var request []byte

	if criteria != nil {
		request = []byte(criteria.(string))
	}

	response, err = http.Post("http://localhost:2323", "application/json",
		bytes.NewReader(request))
	return err
}

func iShouldReceiveABadRequestMessage() error {
	if response.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("Should have received a bad response")
	}
	return nil
}

func iHaveAValidSearchCriteria() error {
	criteria = `{"query": "Fat Freddy's Cat"}`
	return nil
}

func iShouldReceiveAListOfKittens() error {
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Should have received a status OK %v",
			response.StatusCode)
	}

	var body handlers.SearchResponse

	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&body)
	if err != nil || len(body.Kittens) < 1 {
		return fmt.Errorf("Should have received a list of kittens: %v", err)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have no search criteria$`, iHaveNoSearchCriteria)
	s.Step(`^I call the search endpoint$`, iCallTheSearchEndpoint)
	s.Step(`^I should receive a bad request message$`, iShouldReceiveABadRequestMessage)
	s.Step(`^I have a valid search criteria$`, iHaveAValidSearchCriteria)
	s.Step(`^I should receive a list of kittens$`, iShouldReceiveAListOfKittens)

	s.BeforeScenario(func(interface{}) {
		cleanDB()
		err := setupDB()
		if err != nil {
			log.Fatal(err)
		}

		startServer()
		fmt.Printf("Server running with pid: %v", server.Process.Pid)
	})

	s.AfterScenario(func(interface{}, error) {
		server.Process.Kill()
	})

	waitForDB()
}

var server *exec.Cmd
var store *data.MongoStore

func startServer() {
	server = exec.Command("go", "run", "../main.go")

	go server.Run()
	time.Sleep(3 * time.Second)
}

func cleanDB() error {
	return store.DeleteAllKittens()
}

func setupDB() error {
	kittens := []data.Kitten{
		data.Kitten{
			ID:     "1",
			Name:   "Felix",
			Weight: 12.3,
		},
		data.Kitten{
			ID:     "2",
			Name:   "Fat Freddy's Cat",
			Weight: 20.0,
		},
		data.Kitten{
			ID:     "3",
			Name:   "Garfield",
			Weight: 35.0,
		},
	}

	return store.InsertKittens(kittens...)
}

func waitForDB() {
	var err error

	for i := 0; i < 10; i++ {
		store, err = data.NewMongoStore("mongodb://localhost:27017/kittenserver")
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
