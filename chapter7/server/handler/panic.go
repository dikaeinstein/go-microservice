package handler

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/alexcesaro/statsd"
	"github.com/sirupsen/logrus"
)

type panicHandler struct {
	next   http.Handler
	statsd *statsd.Client
	logger *logrus.Logger
}

// NewPanic creates a panic handler with the given statsd client, logger and next handler
func NewPanic(statsd *statsd.Client, logger *logrus.Logger, next http.Handler) http.Handler {
	return &panicHandler{next, statsd, logger}
}

func (p *panicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			p.logger.WithFields(logrus.Fields{
				"handler": "Panic",
				"status":  http.StatusInternalServerError,
				"method":  r.Method,
				"path":    r.URL.Path,
				"query":   r.URL.RawQuery,
			}).Error(fmt.Sprintf("Error: %v\n%s", err, debug.Stack()))

			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

	p.next.ServeHTTP(w, r)
}
