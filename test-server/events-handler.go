package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sammhicks/network-broadcast/messages"
)

type eventsHandler struct {
	pub *messages.Publisher
}

func (eh *eventsHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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

	sub, ok := eh.pub.NewSubscriber()

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
}
