package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/sammhicks/network-broadcast/config"
	"github.com/sammhicks/network-broadcast/messages"
)

type serverConfig struct {
	Port int64
}

func main() {
	var conf serverConfig

	err := config.Load(&conf)

	if err != nil {
		log.Fatalln("Error reading config: ", err)
	}

	var allDone sync.WaitGroup

	mc := make(chan messages.Message)

	pub := messages.PublishMessages(mc, &allDone)

	httpHandler := http.NewServeMux()

	httpHandler.Handle("/message", &messageHandler{mc: mc})

	httpHandler.Handle("/events", &eventsHandler{pub: pub})

	httpHandler.Handle("/", http.FileServer(http.Dir("static")))

	s := &http.Server{
		Addr:    net.JoinHostPort("", strconv.FormatInt(conf.Port, 10)),
		Handler: httpHandler,
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
		signal.Stop(signals)

		log.Println("Shutting down server")

		s.Shutdown(context.Background())
	}()

	s.ListenAndServe()

	close(mc)

	log.Println("Waiting for all routines to terminate")

	allDone.Wait()
}
