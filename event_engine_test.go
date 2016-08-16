package main

import "testing"

func string_arrays_equal(lhs []string, rhs []string) bool {
    if len(lhs) != len(rhs) { return false }
    for i := range lhs {
        if lhs[i] != rhs[i] {
            return false
        }
    }

    return true
}

func TestEventEngineStep(t *testing.T) {
    engine := EventEngine{}

    handlerTriggered := false
    testHandler := func (EventTrigger, string, *EventData) {
        handlerTriggered = true
    }

    engine.RegisterHandler("foo.bar", testHandler)
    engine.Trigger("foo.bar", nil)

    if !string_arrays_equal(engine.QueuedEvents(), []string{"foo.bar"}) {
        t.Errorf("expected event queue to be %v, got %v", []string{"foo.bar"}, engine.QueuedEvents())
    }

    numEventsProcessed := engine.Step()

    if numEventsProcessed != 1 {
        t.Errorf("expected %v events to be processed, got %v", 1, numEventsProcessed)
    }

    if !handlerTriggered {
        t.Error("expected handler to be triggered")
    }
}
