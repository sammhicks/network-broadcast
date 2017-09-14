package messages

import (
	"log"
	"sync"
)

// Subscriber subscribes to a Publisher, and has a channel that receives messages
type Subscriber interface {
	Messages() <-chan Message
	Unsubscribe()
	Done()
}

type subscriber struct {
	messages    chan Message
	unSubscribe chan *subscriber
	done        *sync.WaitGroup
}

// MessageChannel returns the channel in which the publisher publishes messages
func (sub *subscriber) Messages() <-chan Message {
	return sub.messages
}

// Unsubscribe removes this subscriber from its publisher's list
func (sub *subscriber) Unsubscribe() {
	for {
		select {
		case sub.unSubscribe <- sub:
			return
		case _, ok := <-sub.messages:
			if !ok {
				return
			}
		}
	}
}

func (sub *subscriber) Done() {
	log.Println("Sub Done")
	sub.done.Done()
}

// Publisher publishes messages to subscribers
type Publisher interface {
	NewSubscriber() (Subscriber, bool)

	Done() <-chan struct{}
}

type publisher struct {
	done         <-chan struct{}
	newSub       chan *subscriber
	unSubRequest chan *subscriber
	allDone      *sync.WaitGroup
}

// NewSubscriber creates a new Subscriber which is subscribed to the Publisher
func (p *publisher) NewSubscriber() (Subscriber, bool) {
	s := &subscriber{
		messages:    make(chan Message),
		unSubscribe: p.unSubRequest,
		done:        p.allDone}

	select {
	case <-p.done:
		return nil, false
	case p.newSub <- s:
		return s, true
	}
}

func (p *publisher) Done() <-chan struct{} {
	return p.done
}

// PublishMessages starts published messages received from messageSrc
func PublishMessages(messageSrc <-chan Message, pubDone *sync.WaitGroup) Publisher {
	pubDone.Add(1)
	subs := make(map[Subscriber]chan<- Message)

	done := make(chan struct{})
	newSub := make(chan *subscriber)
	unSubRequest := make(chan *subscriber)

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
			case sub := <-unSubRequest:
				log.Println("Dead Subscriber")
				close(sub.messages)
				delete(subs, sub)
			}
		}
	}()

	return &publisher{
		done:         done,
		newSub:       newSub,
		unSubRequest: unSubRequest,
		allDone:      pubDone}
}
