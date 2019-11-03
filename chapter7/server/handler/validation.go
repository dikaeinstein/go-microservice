package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/dikaeinstein/go-microservice/chapter7/server/serialize"

	"github.com/dikaeinstein/go-microservice/chapter7/server/entity"

	"github.com/alexcesaro/statsd"
	"github.com/sirupsen/logrus"
)

type validation struct {
	next   http.Handler
	statsd *statsd.Client
	logger *logrus.Logger
}

// NewValidation creates a validation handler with the given statsd client, logger and next handler
func NewValidation(statsd *statsd.Client, logger *logrus.Logger, next http.Handler) http.Handler {
	return &validation{next, statsd, logger}
}

// Name is the context type for name
type Name string

func (v validation) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request entity.HelloWorldRequest

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		v.statsd.Increment(validationFailed)

		message := serialize.NewSerializableRequest(r)
		v.logger.WithFields(logrus.Fields{
			"handler": "Validation",
			"status":  http.StatusBadRequest,
			"method":  r.Method,
		}).Info(message.ToJSON())
	}

	c := context.WithValue(r.Context(), Name("name"), request.Name)
	r = r.WithContext(c)

	v.statsd.Increment(validationSuccess)
	v.next.ServeHTTP(w, r)
}
