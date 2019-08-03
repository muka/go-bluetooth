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

var GattService1Interface = "org.bluez.GattService1"


// NewGattService1 create a new instance of GattService1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/serviceXX
func NewGattService1(objectPath string) (*GattService1, error) {
	a := new(GattService1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattService1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattService1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// GattService1 Service hierarchy
// GATT remote and local service representation. Object path for local services
// is freely definable.
// External applications implementing local services must register the services
// using GattManager1 registration method and must implement the methods and
// properties defined in GattService1 interface.
type GattService1 struct {
	client     *bluez.Client
	Properties *GattService1Properties
}

// GattService1Properties contains the exposed properties of an interface
type GattService1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Characteristics 
	Characteristics []dbus.ObjectPath `dbus:"emit"`

	// IsService 
	IsService bool `dbus:"ignore"`

	// UUID 128-bit service UUID.
	UUID string

	// Primary Indicates whether or not this GATT service is a
  // primary service. If false, the service is secondary.
	Primary bool

	// Device Object path of the Bluetooth device the service
  // belongs to. Only present on services from remote
  // devices.
	Device []dbus.ObjectPath `dbus:"ignore=isService"`

	// Includes Array of object paths representing the included
  // services of this service.
	Includes []dbus.ObjectPath

}

func (p *GattService1Properties) Lock() {
	p.lock.Lock()
}

func (p *GattService1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *GattService1) Close() {
	a.client.Disconnect()
}


// ToMap convert a GattService1Properties to map
func (a *GattService1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an GattService1Properties
func (a *GattService1Properties) FromMap(props map[string]interface{}) (*GattService1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattService1Properties
func (a *GattService1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattService1Properties, error) {
	s := new(GattService1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *GattService1) GetProperties() (*GattService1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattService1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattService1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *GattService1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *GattService1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



