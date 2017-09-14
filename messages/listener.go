package messages

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

func runSubscriber(conn net.Conn, sub Subscriber) {
	defer sub.Done()
	defer conn.Close()

	bufferedMessages := bufferMessage(sub.Messages())

	encoder := json.NewEncoder(conn)

	for m := range bufferedMessages {
		err := encoder.Encode(m)

		if err != nil {
			sub.Unsubscribe()
			return
		}
	}
}

// ListenForSubscribers will attach a subscriber to each connection to the listener
func ListenForSubscribers(connectionListener net.Listener, pub Publisher, finished *sync.WaitGroup) {
	// Add the listener and listener closer
	finished.Add(2)
	defer finished.Done()
	defer log.Println("Listener Done")

	var isAlive atomic.Value
	isAlive.Store(true)

	// Run the listener closer in the background
	go func() {
		defer finished.Done()
		defer log.Println("Listener Closer")

		<-pub.Done()

		isAlive.Store(false)

		err := connectionListener.Close()

		if err != nil {
			log.Printf("Error closing Listener %v: \"%v\"\n", connectionListener, err)
		}
	}()

	for {
		conn, err := connectionListener.Accept()

		if err != nil {
			if isAlive.Load().(bool) {
				log.Printf("Listener \"%v\" failed unexpectedly\n", connectionListener)
			}
			break
		}

		log.Println("New Connection")

		sub, ok := pub.NewSubscriber()

		if ok {
			go runSubscriber(conn, sub)
		}
	}
}
