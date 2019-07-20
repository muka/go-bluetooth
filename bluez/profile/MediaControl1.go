package profile

import (
	"github.com/muka/go-bluetooth/bluez"
)

// NewMediaControl1 create a new MediaControl1 client
func NewMediaControl1(hostID string) *MediaControl1 {
	a := new(MediaControl1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: bluez.MediaControl1Interface,
			Path:  "/org/bluez/" + hostID,
			Bus:   bluez.SystemBus,
		},
	)
	a.Properties = new(MediaControl1Properties)
	a.GetProperties()
	return a
}

// MediaControl1 client
type MediaControl1 struct {
	client     *bluez.Client
	Properties *MediaControl1Properties
}

//MediaControl1Properties contains the exposed properties of an interface
type MediaControl1Properties struct {
	Connected bool
	Player    string
}

// Close the connection
func (a *MediaControl1) Close() {
	a.client.Disconnect()
}

//GetProperties load all available properties
func (a *MediaControl1) GetProperties() (*MediaControl1Properties, error) {
	err := a.client.GetProperties(a.Properties)
	return a.Properties, err
}

//SetProperty set a property
func (a *MediaControl1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}
