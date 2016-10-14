package bluez

import (
	"github.com/muka/device-manager/client"
	"log"
)

// NewAdapter1 create a new Adapter1 client
func NewAdapter1(idx int) *Adapter1 {
	a := new(Adapter1)
	a.client = client.NewClient("org.bluez.Adapter1", "/org/bluez/hci"+string(idx))
	return a
}

// Adapter1 client
type Adapter1 struct {
	client *client.Client
	logger *log.Logger
}

//StartDiscovery on the adapter
func (a *Adapter1) StartDiscovery() error {
	err := a.client.Call("StartDiscovery", 0).Store()
	return err
}

//StopDiscovery on the adapter
func (a *Adapter1) StopDiscovery() error {
	err := a.client.Call("StartDiscovery", 0).Store()
	return err
}
