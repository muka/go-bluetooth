package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/util"
)

var log = util.NewLogger("main")

func main() {
	//
	// emitter.On("device", func(ev emitter.Event) {
	// 	log.Printf("Got event %s: %s", ev.Name, ev.Data)
	// })
	//
	// err := api.StartDiscovery()
	// if err != nil {
	// 	panic(err)
	// }

	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='/org/freedesktop/DBus',interface='org.freedesktop.DBus',sender='org.freedesktop.DBus'")

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		fmt.Println(v)
	}

}
