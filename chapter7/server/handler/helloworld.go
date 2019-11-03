package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/alexcesaro/statsd"
	"github.com/dikaeinstein/go-microservice/chapter7/server/entity"
	"github.com/dikaeinstein/go-microservice/chapter7/server/serialize"
	"github.com/sirupsen/logrus"
)

type helloWorld struct {
	statsd *statsd.Client
	logger *logrus.Logger
}

// NewHelloWorld creates a new handler with the given statsd client and logger
func NewHelloWorld(statsd *statsd.Client, logger *logrus.Logger) http.Handler {
	return &helloWorld{statsd, logger}
}

func (h *helloWorld) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	timing := h.statsd.NewTiming()

	name := r.Context().Value(Name("name")).(string)
	response := entity.HelloWorldResponse{Message: "Hello " + name}

	encoder := json.NewEncoder(rw)
	encoder.Encode(response)

	h.statsd.Increment(helloworldSuccess)

	message := serialize.NewSerializableRequest(r)
	h.logger.WithFields(logrus.Fields{
		"handler": "HelloWorld",
		"status":  http.StatusOK,
		"method":  r.Method,
	}).Info(message.ToJSON())

	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

	timing.Send(helloworldTiming)
}
