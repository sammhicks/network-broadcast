package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/sammhicks/network-broadcast/config"
	"github.com/sammhicks/network-broadcast/messages"
)

type clientConfig struct {
	Network string
	Address string
}

func main() {
	var conf clientConfig

	err := config.Load(&conf)

	if err != nil {
		log.Fatalln("Error reading config: ", err)
	}

	conn, err := net.Dial(conf.Network, conf.Address)

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
