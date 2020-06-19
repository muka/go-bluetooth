package obex

import (
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

// NewObexSession1 create a new ObexSession1 client
func NewObexSession1(path string) *ObexSession1 {
	a := new(ObexSession1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: "org.bluez.obex.Session1",
			Path:  dbus.ObjectPath(path),
			Bus:   bluez.SessionBus,
		},
	)
	a.Properties = new(ObexSession1Properties)
	_, err := a.GetProperties()
	if err != nil {
		log.Warn(err)
	}
	return a
}

// ObexSession1 client
type ObexSession1 struct {
	client     *bluez.Client
	Properties *ObexSession1Properties
}

// ObexSession1Properties exposed properties for ObexSession1
type ObexSession1Properties struct {
	Source      string // [readonly] Bluetooth adapter address
	Destination string // [readonly] Bluetooth device address
	Channel     byte   // [readonly] Bluetooth channel
	Target      string // [readonly] Target UUID
	Root        string // [readonly] Root path
}

// Close the connection
func (d *ObexSession1) Close() {
	d.client.Disconnect()
}

//GetProperties load all available properties
func (d *ObexSession1) GetProperties() (*ObexSession1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}

//GetProperty get a property
func (d *ObexSession1) GetProperty(name string) (dbus.Variant, error) {
	return d.client.GetProperty(name)
}
