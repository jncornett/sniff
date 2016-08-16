package main

type EventName string

type EventData struct {}

type EventEngine struct {}

type EventTrigger interface {
    Trigger(EventName, *EventData)
}

type EventHandler func(EventTrigger, EventName, *EventData)

func (self EventEngine) RegisterHandler(_ EventName, _ EventHandler) {}

func (self EventEngine) Trigger(_ EventName, _ *EventData) {}

func (self EventEngine) QueuedEvents() []string { return []string{} }

func (self EventEngine) Step() int { return 0 }
