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

package network



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var NetworkServer1Interface = "org.bluez.NetworkServer1"


// NewNetworkServer1 create a new instance of NetworkServer1
//
// Args:
// 	objectPath: /org/bluez/{hci0,hci1,...}
func NewNetworkServer1(objectPath string) (*NetworkServer1, error) {
	a := new(NetworkServer1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: NetworkServer1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(NetworkServer1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// NetworkServer1 Network server hierarchy

type NetworkServer1 struct {
	client     *bluez.Client
	Properties *NetworkServer1Properties
}

// NetworkServer1Properties contains the exposed properties of an interface
type NetworkServer1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *NetworkServer1Properties) Lock() {
	p.lock.Lock()
}

func (p *NetworkServer1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *NetworkServer1) Close() {
	a.client.Disconnect()
}


// ToMap convert a NetworkServer1Properties to map
func (a *NetworkServer1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an NetworkServer1Properties
func (a *NetworkServer1Properties) FromMap(props map[string]interface{}) (*NetworkServer1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an NetworkServer1Properties
func (a *NetworkServer1Properties) FromDBusMap(props map[string]dbus.Variant) (*NetworkServer1Properties, error) {
	s := new(NetworkServer1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *NetworkServer1) GetProperties() (*NetworkServer1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *NetworkServer1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *NetworkServer1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *NetworkServer1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *NetworkServer1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Register Register server for the provided UUID. Every new
// connection to this server will be added the bridge
// interface.
// Valid UUIDs are "gn", "panu" or "nap".
// Initially no network server SDP is provided. Only
// after this method a SDP record will be available
// and the BNEP server will be ready for incoming
// connections.
func (a *NetworkServer1) Register(uuid string, bridge string) error {
	
	return a.client.Call("Register", 0, uuid, bridge).Store()
	
}

//Unregister Unregister the server for provided UUID.
// All servers will be automatically unregistered when
// the calling application terminates.
func (a *NetworkServer1) Unregister(uuid string) error {
	
	return a.client.Call("Unregister", 0, uuid).Store()
	
}

