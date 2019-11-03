package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dikaeinstein/go-microservice/chapter7/server/handler"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"

	"github.com/alexcesaro/statsd"
	"github.com/sirupsen/logrus"
)

func main() {
	const PORT = 8091
	statsd, err := createStatsDClient(os.Getenv("STATSD"))
	if err != nil {
		log.Fatal("Unable to create statsd client")
	}

	logger, err := createLogger(os.Getenv("LOGSTASH"))
	if err != nil {
		log.Fatal("Unable to create logstash client")
	}

	setupHandlers(statsd, logger)

	log.Printf("Server starting on port %v\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil))
}

func setupHandlers(statsd *statsd.Client, logger *logrus.Logger) {
	validation := handler.NewValidation(
		statsd, logger,
		handler.NewHelloWorld(statsd, logger),
	)
	panicHandler := handler.NewPanic(
		statsd, logger,
		handler.NewBangHandler(),
	)

	http.Handle("/helloworld", handler.NewTagRequest(validation))
	http.Handle("/bang", handler.NewTagRequest(panicHandler))
}

func createStatsDClient(addr string) (*statsd.Client, error) {
	return statsd.New(statsd.Address(addr))
}

func createLogger(addr string) (*logrus.Logger, error) {
	retryCount := 0

	l := logrus.New()
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	for ; retryCount < 50; retryCount++ {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			hook := logrustash.New(
				conn,
				logrustash.DefaultFormatter(logrus.Fields{
					"hostname": hostname,
				}),
			)
			l.Hooks.Add(hook)
			return l, err
		}
		log.Println("Unable to connect to logstash, retrying")
		time.Sleep(1 * time.Second)
	}

	log.Fatal("Unable to connect to logstash")
	return nil, err
}
