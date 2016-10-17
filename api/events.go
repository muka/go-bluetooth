package api

import (
	"github.com/muka/bluez-client/util"
)

var log = util.NewLogger("events")

//Callback is a function to be invoked when an event happens
type Callback func(ev Event)

//Event contains information about what happened
type Event struct {
	Name string
	Data interface{}
}

var pipe chan Event
var events = make(map[string][]Callback, 0)

func loop() {
	//log.Printf("loop: Started\n")
	for {
		if pipe == nil {
			//log.Printf("loop: Closed\n")
			return
		}
		ev := <-pipe
		//log.Printf("loop: Trigger event %s\n", ev.Name)
		if events[ev.Name] != nil {
			size := len(events[ev.Name])
			if size == 0 {
				return
			}
			for i := 0; i < size; i++ {
				//log.Printf("Call fn\n")
				events[ev.Name][i](ev)
			}
		}
	}
}

func getPipe() {
	if pipe == nil {
		pipe = make(chan Event)
		go loop()
	}
}

//On registers to an event
func On(event string, callback Callback) {

	if event == "" {
		panic("Cannot use an empty string as event name")
	}

	if events[event] == nil {
		getPipe()
		events[event] = make([]Callback, 0)
	}

	events[event] = append(events[event], callback)
	//log.Printf("Added %s events, len is %d\n", event, len(events[event]))
}

// Emit an event
func Emit(event string, data interface{}) {
	//log.Printf("Emit event %s\n", event)
	getPipe()
	ev := Event{event, data}
	pipe <- ev
}

//Off Removes all callbacks from an event
func Off(name string) {
	//log.Printf("Off %s", name)
	if name == "*" {
		for name := range events {
			if name != "*" {
				Off(name)
			}
		}
	}

	if events[name] != nil {
		delete(events, name)
	}

	if len(events) == 0 {
		close(pipe)
		pipe = nil // will stop the go routine
	}

}
