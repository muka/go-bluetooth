package bluez

import (
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/client"
	"log"
)

// NewAdapter1 create a new Adapter1 client
func NewAdapter1(idx string) *Adapter1 {
	a := new(Adapter1)
	a.client = client.NewClient(
		&client.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.Adapter1",
			Path:  "/org/bluez/hci" + idx,
			Bus:   client.SystemBus,
		},
	)
	a.client.Connect()
	return a
}

// Adapter1 client
type Adapter1 struct {
	client *client.Client
	logger *log.Logger
}

//StartDiscovery on the adapter
func (a *Adapter1) StartDiscovery() error {
	return a.client.Call("StartDiscovery", 0).Store()
}

//StopDiscovery on the adapter
func (a *Adapter1) StopDiscovery() error {
	return a.client.Call("StopDiscovery", 0).Store()
}

//RemoveDevice from the list
func (a *Adapter1) RemoveDevice(device string) error {
	return a.client.Call("RemoveDevice", 0, dbus.ObjectPath(device)).Store()
}
