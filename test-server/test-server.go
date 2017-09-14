package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/sammhicks/network-broadcast/messages"
)

func main() {
	var allDone sync.WaitGroup

	mc := make(chan messages.Message)

	pub := messages.PublishMessages(mc, &allDone)

	addr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8080")

	if err != nil {
		log.Fatalln("Cannot resolve address")
	}

	ln, err := net.ListenTCP("tcp4", addr)

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
