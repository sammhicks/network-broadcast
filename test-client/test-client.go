package main

import (
	"encoding/json"
	"log"

	"github.com/sammhicks/network-broadcast/config"
	"github.com/sammhicks/network-broadcast/messages/http"
)

type clientConfig struct {
	URL string
}

func main() {
	var conf clientConfig

	err := config.Load(&conf)

	if err != nil {
		log.Fatalln("Error reading config: ", err)
	}

	done := make(chan struct{})

	for m := range http.ConnectToPublisher(conf.URL, done) {
		b, err := json.Marshal(m)

		if err == nil {
			log.Println("Message: ", string(b))
		} else {
			log.Println("Error: ", err)
		}
	}
}
