package http

import (
	"encoding/json"
	"net/http"

	"github.com/sammhicks/network-broadcast/messages"
	sse "github.com/sammhicks/server-sent-events"
)

type publisher struct {
	pub messages.Publisher
}

func (pub *publisher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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

	er := sse.NewWriter(rw)

	sub, ok := pub.pub.NewSubscriber()

	if !ok {
		http.Error(rw, "Server is shutting down", http.StatusServiceUnavailable)
		return
	}

	defer sub.Done()

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	for {
		select {
		case m, ok := <-sub.Messages():
			data, _ := json.Marshal(m)
			if ok {
				er.WriteEvent(&sse.Event{Event: "message", Data: string(data)})

				f.Flush()
			} else {
				return
			}
		case <-cn.CloseNotify():
			return
		}
	}
}

// HandleSubscriber creates a Http Handler from a publisher that returns an event stream of messages
func HandleSubscriber(pub messages.Publisher) http.Handler {
	return &publisher{pub}
}
