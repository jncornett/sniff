package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

const (
	EQ_BUFFER_SIZE         = 24
	ENDPOINT_PUSH_LISTENER = "/push"
	MY_ADDRESS             = ":8080"
)

type Address string
type Event int
type EventQueue chan Event
type Subscribers map[Address]bool

func (self Address) ToString() string {
	return string(self)
}

type Engine struct {
	eq          EventQueue
	subscribers Subscribers
}

func NewEngine() *Engine {
	return &Engine{make(EventQueue, EQ_BUFFER_SIZE), make(Subscribers)}
}

func (self Engine) Close() {
	close(self.eq)
}

func (self Engine) MainPublisherLoop() {
	for {
		e, more := <-self.eq
		if !more {
			break // channel closed
		}
		// FIXME thread safety of iterating over map that may be mutated?
		for subscriber, subscribed := range self.subscribers {
			if subscribed {
				go func() {
					client, err := rpc.DialHTTP("tcp", subscriber.ToString())
					if err != nil {
						log.Printf("dialing %v: %v", subscriber, err)
						// TODO possibly remove subscriber from map?
						// for now, the event is just 'dropped' as far
						// as this specific subscriber is concerned
						return
					}
					var reply bool
					err = client.Call("Subscriber.Push", e, &reply)
					if err != nil {
						log.Printf("dialing %v: %v", subscriber, err)
						return
					}
					if !reply {
						log.Printf("dialing %v, got result %v", subscriber, reply)
					}
				}()
			}
		}
	}
}

type Service interface {
	Push(Event)
	Subscribe(Address)
	Unsubscribe(Address)
}

func (self Engine) Push(e Event) {
	self.eq <- e
}

func (self Engine) Subscribe(a Address) {
	self.subscribers[a] = true
}

func (self Engine) Unsubscribe(a Address) {
	// FIXME this is not ideal as the time complexity of
	// publishing any event will grow linearly with the
	// total number of subscribers since the birth
	// of the process. This *may* be a partial solution
	// to the problem of iterating over the subscribers
	// that may be mutated by Subscribe/Unsubscribe
	self.subscribers[a] = false
}

func main() {
	engine := NewEngine()
	defer engine.Close()

	// wrap in Service interface for export via RPC.
	// (we wouldn't want to expose Engine.Close(), for example
	var service Service = engine
	rpc.Register(service)
	rpc.HandleHTTP()

	go engine.MainPublisherLoop()
	ln, err := net.Listen("tcp", MY_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	go http.Serve(ln, nil)
}
