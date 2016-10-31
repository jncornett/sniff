package eventing

const (
	EQ_BUFFER_SIZE = 24
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
	WrappedEngine *Engine
}

func (self *Service) Push(e *Event, reply *bool) error {
	self.WrappedEngine.Push(*e)
	*reply = true
	return nil
}

func (self *Service) Subscribe(addr *Address, reply *bool) error {
	self.WrappedEngine.Subscribe(*addr)
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
	self.WrappedEngine.Unsubscribe(*addr)
	*reply = true
	return nil
}
