package main

import "testing"

func TestEventEngineStep(t *testing.T) {
    engine := EventEngine{}

    handlerTriggered := false
    testHandler := func (_ EventEngine) { handlerTriggered = true }

    handle := engine.RegisterHandler("foo.bar", testHandler)

    if handle == nil {
        t.Error("engine could not register handler")
    }

    engine.TriggerEvent("foo.bar")

    if engine.QueuedEvents() != []string{"foo.bar"} {
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

func TestEventEngineStepWhenEmpty(t *testing.T) {
    engine := EventEngine{}
    engine.Step()
}
