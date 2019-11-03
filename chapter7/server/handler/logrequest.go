package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/dikaeinstein/go-microservice/chapter7/server/serialize"

	"github.com/alexcesaro/statsd"
	"github.com/sirupsen/logrus"
)

type logRequest struct {
	statsd *statsd.Client
	logger *logrus.Logger
	next   http.Handler
}

func (l logRequest) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		req := serialize.NewSerializableRequest(r)
		metadata := req.ToMap()

		message := fmt.Sprintf("%s %s %d - %dms",
			r.RemoteAddr, r.Method, r.Response.StatusCode, time.Now().Sub(start))

		if err := recover(); err != nil {
			metadata["error"] = string(debug.Stack())
			l.logger.WithFields(metadata).Error(message)
			rw.WriteHeader(http.StatusInternalServerError)
		}

		l.logger.WithFields(metadata).Info(message)
	}()

	l.next.ServeHTTP(rw, r)
}
