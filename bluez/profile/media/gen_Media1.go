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
)

var Media1Interface = "org.bluez.Media1"


// NewMedia1 create a new instance of Media1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}
func NewMedia1(objectPath string) (*Media1, error) {
	a := new(Media1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Media1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Media1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Media1 Media hierarchy

type Media1 struct {
	client     *bluez.Client
	Properties *Media1Properties
}

// Media1Properties contains the exposed properties of an interface
type Media1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *Media1Properties) Lock() {
	p.lock.Lock()
}

func (p *Media1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Media1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Media1Properties to map
func (a *Media1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Media1Properties
func (a *Media1Properties) FromMap(props map[string]interface{}) (*Media1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Media1Properties
func (a *Media1Properties) FromDBusMap(props map[string]dbus.Variant) (*Media1Properties, error) {
	s := new(Media1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Media1) GetProperties() (*Media1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Media1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Media1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Media1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Media1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//RegisterEndpoint Register a local end point to sender, the sender can
// register as many end points as it likes.
// Note: If the sender disconnects the end points are
// automatically unregistered.
// possible properties:
// string UUID:
// UUID of the profile which the endpoint
// is for.
// byte Codec:
// Assigned number of codec that the
// endpoint implements. The values should
// match the profile specification which
// is indicated by the UUID.
// array{byte} Capabilities:
// Capabilities blob, it is used as it is
// so the size and byte order must match.
// Possible Errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.NotSupported - emitted
// when interface for the end-point is
// disabled.
func (a *Media1) RegisterEndpoint(endpoint dbus.ObjectPath, properties map[string]interface{}) error {
	
	return a.client.Call("RegisterEndpoint", 0, endpoint, properties).Store()
	
}

//UnregisterEndpoint Unregister sender end point.
func (a *Media1) UnregisterEndpoint(endpoint dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterEndpoint", 0, endpoint).Store()
	
}

//RegisterPlayer Register a media player object to sender, the sender
// can register as many objects as it likes.
// Object must implement at least
// org.mpris.MediaPlayer2.Player as defined in MPRIS 2.2
// spec:
// http://specifications.freedesktop.org/mpris-spec/latest/
// Note: If the sender disconnects its objects are
// automatically unregistered.
// Possible Errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.NotSupported
func (a *Media1) RegisterPlayer(player dbus.ObjectPath, properties map[string]interface{}) error {
	
	return a.client.Call("RegisterPlayer", 0, player, properties).Store()
	
}

//UnregisterPlayer Unregister sender media player.
func (a *Media1) UnregisterPlayer(player dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterPlayer", 0, player).Store()
	
}

