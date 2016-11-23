package emitter

import "github.com/tj/go-debug"

var dbg = debug.Debug("bluez:emitter")

//Callback is a function to be invoked when an event happens
type Callback func(ev Event)

//Event contains information about what happened
type Event interface {
	GetName() string
	GetData() interface{}
}

//BaseEvent contains unspecialized information about what happened
type BaseEvent struct {
	name string
	data interface{}
}

//GetName return the event name
func (e BaseEvent) GetName() string {
	return e.name
}

//GetData return the event data
func (e BaseEvent) GetData() interface{} {
	return e.data
}

var pipe chan Event
var events = make(map[string][]*Callback, 0)

func loop() {
	dbg("loop: Started")
	for {

		if pipe == nil {
			dbg("loop: Closed")
			return
		}

		dbg("loop: Waiting for events")
		ev := <-pipe
		dbg("loop: Trigger event `%s`", ev.GetName())

		if events[ev.GetName()] != nil {
			size := len(events[ev.GetName()])
			if size == 0 {
				dbg("loop: No callback(s)")
			} else {
				dbg("loop: %d callback(s)", size)
				for i := 0; i < size; i++ {
					cb := *events[ev.GetName()][i]
					go cb(ev)
				}
			}
		}
		dbg("loop: done event trigger")
	}
}

func getPipe() {
	if pipe == nil {
		dbg("Init pipe")
		pipe = make(chan Event, 1)
		go loop()
	}
}

// NewCallback creates a new Callback to be passed to the emitter
func NewCallback(fn func(ev Event)) *Callback {
	cb := Callback(fn)
	return &cb
}

//On registers to an event
func On(event string, callback *Callback) {

	if event == "" {
		panic("Cannot use an empty string as event name")
	}

	if events[event] == nil {
		getPipe()
		events[event] = make([]*Callback, 0)
	}

	events[event] = append(events[event], callback)
	dbg("Added to `%s` event, len is %d", event, len(events[event]))
}

// Emit an event
func Emit(name string, data interface{}) {
	dbg("Emit event `%s` -> %v", name, data)
	getPipe()
	ev := BaseEvent{name, data}
	dbg("Send to pipe")
	pipe <- ev
}

//Off Removes all callbacks from an event
func Off(name string, callback *Callback) {

	dbg("Off %s", name)

	if name == "*" {
		for name := range events {
			if name != "*" {
				Off(name, nil)
			}
		}
	}

	if callback == nil {
		delete(events, name)
	}

	if _, ok := events[name]; ok {
		for i, cb := range events[name] {
			// compare pointers to see if the exactly same function
			if cb == callback {
				events[name] = append(events[name][:i], events[name][i+1:]...)
			}
		}
	}

	if len(events) == 0 {
		close(pipe)
		pipe = nil // will stop the go routine
	}

}
