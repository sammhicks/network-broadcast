package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sammhicks/network-broadcast/messages"
	sse "github.com/sammhicks/server-sent-events"
)

// ConnectToPublisher connects to url and reads events, placing them on the returned channel
func ConnectToPublisher(url string, done <-chan struct{}) <-chan messages.Message {
	out := make(chan messages.Message)

	resp, err := http.Get(url)

	if err != nil {
		close(out)

		return out
	}

	go func() {
		<-done

		resp.Body.Close()
	}()

	go func() {
		defer close(out)

		er := sse.NewReader(resp.Body)

		for {
			ev, err := er.ReadEvent()

			if err != nil {
				return
			}

			var m messages.Message

			err = json.Unmarshal([]byte(ev.Data), &m)

			if err != nil {
				log.Println("Error unmarshalling ", ev)
				return
			}

			out <- m
		}
	}()

	return out
}
