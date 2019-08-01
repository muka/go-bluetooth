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

package sap



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var SimAccess1Interface = "org.bluez.SimAccess1"


// NewSimAccess1 create a new instance of SimAccess1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}
func NewSimAccess1(objectPath string) (*SimAccess1, error) {
	a := new(SimAccess1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: SimAccess1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(SimAccess1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// SimAccess1 Sim Access Profile hierarchy

type SimAccess1 struct {
	client     *bluez.Client
	Properties *SimAccess1Properties
}

// SimAccess1Properties contains the exposed properties of an interface
type SimAccess1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Connected Indicates if SAP client is connected to the server.
	Connected bool

}

func (p *SimAccess1Properties) Lock() {
	p.lock.Lock()
}

func (p *SimAccess1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *SimAccess1) Close() {
	a.client.Disconnect()
}


// ToMap convert a SimAccess1Properties to map
func (a *SimAccess1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an SimAccess1Properties
func (a *SimAccess1Properties) FromMap(props map[string]interface{}) (*SimAccess1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an SimAccess1Properties
func (a *SimAccess1Properties) FromDBusMap(props map[string]dbus.Variant) (*SimAccess1Properties, error) {
	s := new(SimAccess1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *SimAccess1) GetProperties() (*SimAccess1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *SimAccess1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *SimAccess1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *SimAccess1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *SimAccess1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Disconnect Disconnects SAP client from the server.
// Possible errors: org.bluez.Error.Failed
func (a *SimAccess1) Disconnect() error {
	
	return a.client.Call("Disconnect", 0, ).Store()
	
}

