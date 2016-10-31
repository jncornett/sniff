package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/jncornett/sniff/eventing"
)

const (
	MY_ADDRESS = ":8080"
)

func main() {
	engine := eventing.NewEngine()
	defer engine.Close()

	// wrap in Service interface for export via RPC.
	// (we wouldn't want to expose Engine.Close(), for example
	service := eventing.Service{engine}
	rpc.Register(&service)
	rpc.HandleHTTP()

	go eventing.MainPublisherLoop(engine, func(subscriber eventing.Address, e eventing.Event) {
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
	})

	ln, err := net.Listen("tcp", MY_ADDRESS)
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(ln, nil)
}
