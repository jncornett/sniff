package main

import "fmt"

func printDataHandler(_ EventTrigger, eventName string, data *EventData) {
    fmt.Printf("handling event '%v' with data %v", eventName, data)
}

func main() {
    ee := EventEngine{}
    ee.RegisterHandler("foo", printDataHandler)
}
