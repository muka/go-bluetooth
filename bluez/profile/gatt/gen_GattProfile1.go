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

package gatt



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var GattProfile1Interface = "org.bluez.GattProfile1"


// NewGattProfile1 create a new instance of GattProfile1
//
// Args:
// 	servicePath: <application dependent>
// 	objectPath: <application dependent>
func NewGattProfile1(servicePath string, objectPath string) (*GattProfile1, error) {
	a := new(GattProfile1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: GattProfile1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattProfile1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// GattProfile1 GATT Profile hierarchy
// Local profile (GATT client) instance. By registering this type of object
// an application effectively indicates support for a specific GATT profile
// and requests automatic connections to be established to devices
// supporting it.
type GattProfile1 struct {
	client     *bluez.Client
	Properties *GattProfile1Properties
}

// GattProfile1Properties contains the exposed properties of an interface
type GattProfile1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// UUIDs 128-bit GATT service UUIDs to auto connect.
	UUIDs []string

}

func (p *GattProfile1Properties) Lock() {
	p.lock.Lock()
}

func (p *GattProfile1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *GattProfile1) Close() {
	a.client.Disconnect()
}


// ToMap convert a GattProfile1Properties to map
func (a *GattProfile1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an GattProfile1Properties
func (a *GattProfile1Properties) FromMap(props map[string]interface{}) (*GattProfile1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattProfile1Properties
func (a *GattProfile1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattProfile1Properties, error) {
	s := new(GattProfile1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *GattProfile1) GetProperties() (*GattProfile1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattProfile1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattProfile1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *GattProfile1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *GattProfile1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Release This method gets called when the service daemon
// unregisters the profile. The profile can use it to
// do cleanup tasks. There is no need to unregister the
// profile, because when this method gets called it has
// already been unregistered.
func (a *GattProfile1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

