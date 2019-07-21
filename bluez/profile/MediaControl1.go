package profile

import (
	"github.com/muka/go-bluetooth/bluez"
)

/*
Media Control hierarchy
=======================

Service		org.bluez
Interface	org.bluez.MediaControl1
Object path	[variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX

*/

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

// void Play() [Deprecated]
//
//   Resume playback.
func (a *MediaControl1) Play() error {
	return a.client.Call("Play", 0).Store()
}

// void Pause() [Deprecated]
//
//   Pause playback.
func (a *MediaControl1) Pause() error {
	return a.client.Call("Pause", 0).Store()
}

// void Stop() [Deprecated]
//
//   Stop playback.
func (a *MediaControl1) Stop() error {
	return a.client.Call("Stop", 0).Store()
}

// void Next() [Deprecated]
//
//   Next item.
func (a *MediaControl1) Next() error {
	return a.client.Call("Next", 0).Store()
}

// void Previous() [Deprecated]
//
//   Previous item.
func (a *MediaControl1) Previous() error {
	return a.client.Call("Previous", 0).Store()
}

// void VolumeUp() [Deprecated]
//
//   Adjust remote volume one step up
func (a *MediaControl1) VolumeUp() error {
	return a.client.Call("VolumeUp", 0).Store()
}

// void VolumeDown() [Deprecated]
//
//   Adjust remote volume one step down
func (a *MediaControl1) VolumeDown() error {
	return a.client.Call("VolumeDown", 0).Store()
}

// void FastForward() [Deprecated]
//
//   Fast forward playback, this action is only stopped
//   when another method in this interface is called.
func (a *MediaControl1) FastForward() error {
	return a.client.Call("FastForward", 0).Store()
}

// void Rewind() [Deprecated]
//
//   Rewind playback, this action is only stopped
//   when another method in this interface is called.
func (a *MediaControl1) Rewind() error {
	return a.client.Call("Rewind", 0).Store()
}
