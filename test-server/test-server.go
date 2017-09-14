package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sammhicks/network-broadcast/config"
	"github.com/sammhicks/network-broadcast/messages"
)

type serverConfig struct {
	Network string
	Address string
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

	addr, err := net.ResolveTCPAddr(conf.Network, conf.Address)

	if err != nil {
		log.Fatalln("Cannot resolve address")
	}

	ln, err := net.ListenTCP(conf.Network, addr)

	if err != nil {
		log.Fatalln("Cannot open listener")
	}

	go messages.ListenForSubscribers(ln, pub, &allDone)

	for i := 0; i < 60; i++ {
		time.Sleep(500 * time.Millisecond)
		mc <- messages.Message{Sender: "name", Content: fmt.Sprint(i)}
	}

	close(mc)

	log.Println("Shutting Down")

	allDone.Wait()
}
