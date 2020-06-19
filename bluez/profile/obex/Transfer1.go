package obex

import (
	"github.com/godbus/dbus/v5"
	"github.com/muka/go-bluetooth/bluez"
	log "github.com/sirupsen/logrus"
)

// NewObexTransfer1 create a new ObexTransfer1 client
func NewObexTransfer1(path string) *ObexTransfer1 {
	a := new(ObexTransfer1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: "org.bluez.obex.Transfer1",
			Path:  dbus.ObjectPath(path),
			Bus:   bluez.SessionBus,
		},
	)
	a.Properties = new(ObexTransfer1Properties)
	_, err := a.GetProperties()
	if err != nil {
		log.Warn(err)
	}
	return a
}

// ObexTransfer1 client
type ObexTransfer1 struct {
	client     *bluez.Client
	Properties *ObexTransfer1Properties
}

// ObexTransfer1Properties exposed properties for ObexTransfer1
type ObexTransfer1Properties struct {
	Status      string
	Session     dbus.ObjectPath
	Name        string
	Type        string
	Time        uint64
	Size        uint64
	Transferred uint64
	Filename    string
}

// Close the connection
func (d *ObexTransfer1) Close() {
	d.client.Disconnect()
}

//GetProperties load all available properties
func (d *ObexTransfer1) GetProperties() (*ObexTransfer1Properties, error) {
	err := d.client.GetProperties(d.Properties)
	return d.Properties, err
}

//GetProperty get a property
func (d *ObexTransfer1) GetProperty(name string) (dbus.Variant, error) {
	return d.client.GetProperty(name)
}

//
// Stops the current transference.
//
// Possible errors: org.bluez.obex.Error.NotAuthorized
// 	- org.bluez.obex.Error.InProgress
// 	- org.bluez.obex.Error.Failed
//
func (a *ObexTransfer1) Cancel() error {
	return a.client.Call("Cancel", 0).Store()
}

//
// Suspend transference.
//
// Possible errors: org.bluez.obex.Error.NotAuthorized
// org.bluez.obex.Error.NotInProgress
//
// Note that it is not possible to suspend transfers
// which are queued which is why NotInProgress is listed
// as possible error.
//
func (a *ObexTransfer1) Suspend() error {
	return a.client.Call("Suspend", 0).Store()
}

//
// Resume transference.
//
// Possible errors: org.bluez.obex.Error.NotAuthorized
// org.bluez.obex.Error.NotInProgress
//
// Note that it is not possible to resume transfers
// which are queued which is why NotInProgress is listed
// as possible error.
//
func (a *ObexTransfer1) Resume() error {
	return a.client.Call("Resume", 0).Store()
}
