package messages

import (
	"log"
	"sync"
)

// Subscriber subscribes to a Publisher, and has a channel that receives messages
type Subscriber struct {
	messages    chan Message
	unsubscribe chan *Subscriber
	done        *sync.WaitGroup
}

// Messages returns the channel in which the publisher publishes messages
func (sub *Subscriber) Messages() <-chan Message {
	return sub.messages
}

// Unsubscribe removes this Subscriber from its publisher's list
func (sub *Subscriber) Unsubscribe() {
	for {
		select {
		case sub.unsubscribe <- sub:
			return
		case _, ok := <-sub.messages:
			if !ok {
				return
			}
		}
	}
}

// Done signals that this Subscriber goroutine has terminated
func (sub *Subscriber) Done() {
	log.Println("Subscriber Done")
	sub.Unsubscribe()
	sub.done.Done()
}

// Publisher publishes messages to subscribers
type Publisher struct {
	done         <-chan struct{}
	newSub       chan *Subscriber
	unsubRequest chan *Subscriber
	allDone      *sync.WaitGroup
}

// NewSubscriber creates a new Subscriber which is subscribed to the Publisher
func (p *Publisher) NewSubscriber() (*Subscriber, bool) {
	s := &Subscriber{
		messages:    make(chan Message),
		unsubscribe: p.unsubRequest,
		done:        p.allDone}

	select {
	case <-p.done:
		return nil, false
	case p.newSub <- s:
		return s, true
	}
}

// Done returns a channel that is closed when the publisher terminates
func (p *Publisher) Done() <-chan struct{} {
	return p.done
}

// PublishMessages starts published messages received from messageSrc
func PublishMessages(messageSrc <-chan Message, pubDone *sync.WaitGroup) *Publisher {
	pubDone.Add(1)
	subs := make(map[*Subscriber]chan<- Message)

	done := make(chan struct{})
	newSub := make(chan *Subscriber)
	unsubRequest := make(chan *Subscriber)

	go func() {
		defer pubDone.Done()
		defer log.Println("Pub Done")
		for {
			select {
			case newMessage, ok := <-messageSrc:
				if ok {
					for _, messageChannel := range subs {
						messageChannel <- newMessage
					}
				} else {
					close(done)

					for _, messageChannel := range subs {
						close(messageChannel)
					}

					return
				}
			case sub := <-newSub:
				log.Println("New Subscriber")
				pubDone.Add(1)
				subs[sub] = sub.messages
			case sub := <-unsubRequest:
				log.Println("Subscriber Unsubscribing")
				close(sub.messages)
				delete(subs, sub)
			}
		}
	}()

	return &Publisher{
		done:         done,
		newSub:       newSub,
		unsubRequest: unsubRequest,
		allDone:      pubDone}
}
