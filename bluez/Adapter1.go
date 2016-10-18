package bluez

import (
	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/util"
	"github.com/muka/device-manager/client"
	"log"
)

// NewAdapter1 create a new Adapter1 client
func NewAdapter1(hostID string) *Adapter1 {
	a := new(Adapter1)
	a.client = client.NewClient(
		&client.Config{
			Name:  "org.bluez",
			Iface: "org.bluez.Adapter1",
			Path:  "/org/bluez/" + hostID,
			Bus:   client.SystemBus,
		},
	)
	a.Properties = new(Adapter1Properties)
	a.logger = util.NewLogger(hostID)
	return a
}

// Adapter1 client
type Adapter1 struct {
	client     *client.Client
	logger     *log.Logger
	Properties *Adapter1Properties
}

//Adapter1Properties contains the exposed properties of an interface
type Adapter1Properties struct {
	UUIDs               []string
	Discoverable        bool
	Discovering         bool
	Pairable            bool
	Powered             bool
	Address             string
	Alias               string
	Modalias            string
	Name                string
	Class               uint32
	DiscoverableTimeout uint32
	PairableTimeout     uint32
}

// Close the connection
func (a *Adapter1) Close() {
	a.client.Disconnect()
}

//GetProperties load all available properties
func (a *Adapter1) GetProperties() (*Adapter1Properties, error) {
	err := a.client.GetProperties(a.Properties)
	return a.Properties, err
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
