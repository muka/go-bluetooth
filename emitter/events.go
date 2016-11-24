package emitter

import (
	"regexp"
	"strings"
	"sync"

	"github.com/tj/go-debug"
)

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
var mutex = &sync.Mutex{}

func loop() {
	dbg("loop: Started")
	for {

		if pipe == nil {
			dbg("loop: Closed")
			return
		}

		dbg("loop: Waiting for events")
		ev := <-pipe
		if ev == nil {
			dbg("loop: nil event, quit")
			return
		}

		dbg("loop: Trigger event `%s`", ev.GetName())

		mutex.Lock()
		if _, ok := events[ev.GetName()]; ok {
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
		mutex.Unlock()
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

	mutex.Lock()

	if _, ok := events[event]; !ok {
		getPipe()
		events[event] = make([]*Callback, 0)
	}

	events[event] = append(events[event], callback)
	mutex.Unlock()
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

//MatchListeners return a list of matching event names
// replacing * with any char and assuming a namespacing built with dots (.)
// eg. device_name.uu-id-val
func MatchListeners(path string) []string {
	var foundMatches []string
	reg := regexp.MustCompile("^" + strings.Replace(path, "*", ".*", -1) + "$")
	for name := range events {
		if reg.MatchString(name) {
			foundMatches = append(foundMatches, name)
		}
	}
	return foundMatches
}

//RemoveListeners drop a list of listeners by event name
func RemoveListeners(pattern string, callback *Callback) {
	paths := MatchListeners(pattern)
	for _, name := range paths {
		Off(name, callback)
	}
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

	mutex.Lock()

	if callback == nil {
		delete(events, name)
	}

	if _, ok := events[name]; ok {
		for i, cb := range events[name] {
			// compare pointers to see if the exactly same function
			if cb == callback {
				dbg("Drop callback for `%s`", name)
				events[name] = append(events[name][:i], events[name][i+1:]...)
			}
		}
	}

	mutex.Unlock()

	if len(events) == 0 {
		close(pipe)
		pipe = nil // will stop the go routine
	}

}
