// WARNING: generated code, do not edit!
// Copyright © 2019 luca capra
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package media



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
  "fmt"
)

var MediaControl1Interface = "org.bluez.MediaControl1"


// NewMediaControl1 create a new instance of MediaControl1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewMediaControl1(objectPath string) (*MediaControl1, error) {
	a := new(MediaControl1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: MediaControl1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaControl1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}

// NewMediaControl1FromAdapterID create a new instance of MediaControl1
// adapterID: ID of an adapter eg. hci0
func NewMediaControl1FromAdapterID(adapterID string) (*MediaControl1, error) {
	a := new(MediaControl1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: MediaControl1Interface,
			Path:  fmt.Sprintf("/org/bluez/%s", adapterID),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaControl1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// MediaControl1 Media Control hierarchy

type MediaControl1 struct {
	client     *bluez.Client
	Properties *MediaControl1Properties
}

// MediaControl1Properties contains the exposed properties of an interface
type MediaControl1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Connected 
	Connected bool

	// Player Addressed Player object path.
	Player dbus.ObjectPath

}

func (p *MediaControl1Properties) Lock() {
	p.lock.Lock()
}

func (p *MediaControl1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *MediaControl1) Close() {
	a.client.Disconnect()
}


// ToMap convert a MediaControl1Properties to map
func (a *MediaControl1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MediaControl1Properties
func (a *MediaControl1Properties) FromMap(props map[string]interface{}) (*MediaControl1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MediaControl1Properties
func (a *MediaControl1Properties) FromDBusMap(props map[string]dbus.Variant) (*MediaControl1Properties, error) {
	s := new(MediaControl1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MediaControl1) GetProperties() (*MediaControl1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MediaControl1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MediaControl1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *MediaControl1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *MediaControl1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Play Resume playback.
func (a *MediaControl1) Play() error {
	
	return a.client.Call("Play", 0, ).Store()
	
}

//Pause Pause playback.
func (a *MediaControl1) Pause() error {
	
	return a.client.Call("Pause", 0, ).Store()
	
}

//Stop Stop playback.
func (a *MediaControl1) Stop() error {
	
	return a.client.Call("Stop", 0, ).Store()
	
}

//Next Next item.
func (a *MediaControl1) Next() error {
	
	return a.client.Call("Next", 0, ).Store()
	
}

//Previous Previous item.
func (a *MediaControl1) Previous() error {
	
	return a.client.Call("Previous", 0, ).Store()
	
}

//VolumeUp Adjust remote volume one step up
func (a *MediaControl1) VolumeUp() error {
	
	return a.client.Call("VolumeUp", 0, ).Store()
	
}

//VolumeDown Adjust remote volume one step down
func (a *MediaControl1) VolumeDown() error {
	
	return a.client.Call("VolumeDown", 0, ).Store()
	
}

//FastForward Fast forward playback, this action is only stopped
// when another method in this interface is called.
func (a *MediaControl1) FastForward() error {
	
	return a.client.Call("FastForward", 0, ).Store()
	
}

//Rewind Rewind playback, this action is only stopped
// when another method in this interface is called.
// Properties
// boolean Connected [readonly]
// object Player [readonly, optional]
// Addressed Player object path.
func (a *MediaControl1) Rewind() error {
	
	return a.client.Call("Rewind", 0, ).Store()
	
}

