package eventing

const (
	EQ_BUFFER_SIZE = 24
	MY_ADDRESS     = ":8080"
)

type Address string
type Event int
type EventQueue chan Event
type Subscribers map[Address]bool

// Utility to smooth over eventual conversion of
// Address to composite type
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

func (self *Engine) Close() {
	close(self.eq)
}

func (self *Engine) Subscribe(addr Address) {
	self.subscribers[addr] = true
}

func (self *Engine) Unsubscribe(addr Address) {
	self.subscribers[addr] = false
}

func (self *Engine) EachSubscriber(eachFunc func(Address)) {
	for subscriber, subscribed := range self.subscribers {
		if subscribed {
			eachFunc(subscriber)
		}
	}
}

func (self *Engine) Subscribed(addr Address) bool {
	val, ok := self.subscribers[addr]
	return ok && val
}

func (self *Engine) Push(e Event) {
	self.eq <- e
}

func (self *Engine) GetReaderChannel() <-chan Event {
	return self.eq
}

func MainPublisherLoop(engine *Engine, rpcFunc func(Address, Event)) {
	for {
		e, more := <-engine.GetReaderChannel()
		if !more {
			break // channel closed
		}

		// FIXME thread safety of iterating over map that may be mutated?
		engine.EachSubscriber(func(subscriber Address) {
			go rpcFunc(subscriber, e)
		})
	}
}

type Service struct {
	engine *Engine
}

func (self *Service) Push(e *Event, reply *bool) error {
	self.engine.Push(*e)
	*reply = true
	return nil
}

func (self *Service) Subscribe(addr *Address, reply *bool) error {
	self.engine.Subscribe(*addr)
	*reply = true
	return nil
}

func (self *Service) Unsubscribe(addr *Address, reply *bool) error {
	// FIXME this is not ideal as the time complexity of
	// publishing any event will grow linearly with the
	// total number of subscribers since the birth
	// of the process. This *may* be a partial solution
	// to the problem of iterating over the subscribers
	// that may be mutated by Subscribe/Unsubscribe
	self.engine.Unsubscribe(*addr)
	*reply = true
	return nil
}

// func main() {
// 	engine := NewEngine()
// 	defer engine.Close()

// 	// wrap in Service interface for export via RPC.
// 	// (we wouldn't want to expose Engine.Close(), for example
// 	service := Service{engine}
// 	rpc.Register(&service)
// 	rpc.HandleHTTP()

// 	go MainPublisherLoop(engine, func(subscriber Address, e Event) {
// 		client, err := rpc.DialHTTP("tcp", subscriber.ToString())
// 		if err != nil {
// 			log.Printf("dialing %v: %v", subscriber, err)
// 			// TODO possibly remove subscriber from map?
// 			// for now, the event is just 'dropped' as far
// 			// as this specific subscriber is concerned
// 			return
// 		}
// 		var reply bool
// 		err = client.Call("Subscriber.Push", e, &reply)
// 		if err != nil {
// 			log.Printf("dialing %v: %v", subscriber, err)
// 			return
// 		}
// 		if !reply {
// 			log.Printf("dialing %v, got result %v", subscriber, reply)
// 		}
// 	})

// 	ln, err := net.Listen("tcp", MY_ADDRESS)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	http.Serve(ln, nil)
// }
