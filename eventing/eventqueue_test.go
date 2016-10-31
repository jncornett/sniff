package eventing

import "testing"

func TestNewEngineHasNoSubscribers(t *testing.T) {
	engine := NewEngine()
	addr := Address("foo")

	if engine.Subscribed(addr) {
		t.Error("Expecting new engine not to be subscribed at all")
	}
}

func TestNewEngineHasAnOpenChannel(t *testing.T) {
	t.Skip("There is unfortunately no simple (read unhacky) way to check if a channel is closed")
	engine := NewEngine()

	_, open := <-engine.GetReaderChannel()
	if !open {
		t.Error("Expected new engine to have an open channel")
	}
}

func TestEngineClose(t *testing.T) {
	t.Skip("A failed test is costly as we will block (for 10s) until the test panics")
	engine := NewEngine()

	engine.Close()

	_, open := <-engine.GetReaderChannel()
	if open {
		t.Error("Expected engine channel to have been closed")
	}
}

func TestEngineSubcribe(t *testing.T) {
	engine := NewEngine()
	addr := Address("foo")

	engine.Subscribe(addr)

	if !engine.Subscribed(addr) {
		t.Errorf("Expecting engine to be subscribed by %v", addr)
	}
}

func TestEngineUnsubscribe(t *testing.T) {
	engine := NewEngine()
	addr := Address("foo")
	engine.Subscribe(addr)

	engine.Unsubscribe(addr)

	if engine.Subscribed(addr) {
		t.Errorf("Expecting engine to be unsubscribed by %v", addr)
	}
}

func TestEngineSubscribed(t *testing.T) {
	engine := NewEngine()
	addr := Address("foo")

	if engine.Subscribed(addr) {
		t.Errorf("Expecting engine to not be subscribed by %v by default", addr)
	}

	engine.Subscribe(addr)
	if !engine.Subscribed(addr) {
		t.Errorf("Expecting engine to be subscribed by %v", addr)
	}

	engine.Unsubscribe(addr)
	if engine.Subscribed(addr) {
		t.Errorf("Expecting engine to be unsubscribed by %v", addr)
	}
}

func TestEnginePushAndGetReaderChannel(t *testing.T) {
	engine := NewEngine()
	e := Event(13)

	engine.Push(e)

	x, _ := <-engine.GetReaderChannel()
	if x != e {
		t.Errorf("Expected to dequeue %v, not %v", e, x)
	}
}
