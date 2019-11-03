package throttling

import "net/http"

// LimitHandler throttles the number of concurrent requests a handler can handle
type LimitHandler struct {
	connections chan struct{}
	handler     http.Handler
}

// NewLimitHandler constructs an instance of a LimitHandler
func NewLimitHandler(connections int, next http.Handler) *LimitHandler {
	conns := make(chan struct{}, connections)

	for i := 0; i < connections; i++ {
		conns <- struct{}{}
	}

	return &LimitHandler{
		connections: conns,
		handler:     next,
	}
}

func (l *LimitHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	select {
	case <-l.connections:
		l.handler.ServeHTTP(rw, r)
		l.connections <- struct{}{}
	default:
		http.Error(rw, "Busy", http.StatusTooManyRequests)
	}
}
