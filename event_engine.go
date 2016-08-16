package main

// this is a placeholder type for when we actually figure out
// what data to pass to an event handler
type EventData struct {}

type EventEngine struct {}

// View that restricts event handler access to EventEngine
type EventTrigger interface {
    Trigger(string, *EventData)
}

// Event Handler function
type EventHandler func(EventTrigger, string, *EventData)

// EventEngine methods

func (self EventEngine) RegisterHandler(eventName string, _ EventHandler) {}

func (self EventEngine) Trigger(eventName string, data *EventData) {}

func (self EventEngine) QueuedEvents() []string { return []string{} }

func (self EventEngine) Step() int { return 0 }
