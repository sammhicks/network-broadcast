package messages

import (
	"encoding/json"
	"log"
	"net"
)

// ReadFromConn reads messages from conn and places sends them to the returned channel
func ReadFromConn(conn net.Conn) <-chan Message {
	out := make(chan Message)

	go func() {
		defer close(out)
		defer conn.Close()

		dec := json.NewDecoder(conn)

		for dec.More() {
			var m Message

			err := dec.Decode(&m)

			if err != nil {
				log.Println("Error: ", err)
				return
			}

			out <- m
		}
	}()

	return out
}
