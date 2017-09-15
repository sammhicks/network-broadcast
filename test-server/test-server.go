package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/sammhicks/network-broadcast/config"
	"github.com/sammhicks/network-broadcast/messages"
)

type serverConfig struct {
	Port int64
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Welcome to the home page!")
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

	mux := http.NewServeMux()

	mux.HandleFunc("/events", func(rw http.ResponseWriter, req *http.Request) {
		f, ok := rw.(http.Flusher)
		if !ok {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		cn, ok := rw.(http.CloseNotifier)
		if !ok {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		sub, ok := pub.NewSubscriber()

		if !ok {
			http.Error(rw, "Server is shutting down", http.StatusServiceUnavailable)
			return
		}

		defer sub.Done()

		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Connection", "keep-alive")
		rw.WriteHeader(http.StatusOK)

		enc := json.NewEncoder(rw)

		for {
			select {
			case m, ok := <-sub.Messages():
				if ok {
					fmt.Fprintf(rw, "data: ")
					enc.Encode(m)
					fmt.Fprintln(rw, "")
					f.Flush()
				} else {
					return
				}
			case <-cn.CloseNotify():
				return
			}
		}
	})

	mux.Handle("/", http.FileServer(http.Dir("static")))

	go func() {
		for i := 0; ; i++ {
			mc <- messages.Message{Sender: "name", Content: fmt.Sprint(i)}
			time.Sleep(time.Second)
		}
	}()

	s := &http.Server{
		Addr:           ":" + strconv.FormatInt(conf.Port, 10),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
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
