package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/sammhicks/network-broadcast/messages"
)

func main() {
	conn, err := net.Dial("tcp4", "localhost:8080")

	if err != nil {
		log.Fatalln("Could not dial: ", err)
	}

	for m := range messages.ReadFromConn(conn) {
		b, err := json.Marshal(m)

		if err == nil {
			log.Println("Message: ", string(b))
		} else {
			log.Println("Error: ", err)
		}
	}
}
