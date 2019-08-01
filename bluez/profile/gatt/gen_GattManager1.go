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
  "fmt"
)

var GattManager1Interface = "org.bluez.GattManager1"


// NewGattManager1 create a new instance of GattManager1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}
func NewGattManager1(objectPath string) (*GattManager1, error) {
	a := new(GattManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattManager1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}

// NewGattManager1FromAdapterID create a new instance of GattManager1
// adapterID: ID of an adapter eg. hci0
func NewGattManager1FromAdapterID(adapterID string) (*GattManager1, error) {
	a := new(GattManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattManager1Interface,
			Path:  fmt.Sprintf("/org/bluez/%s", adapterID),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// GattManager1 GATT Manager hierarchy
// GATT Manager allows external applications to register GATT services and
// profiles.
// Registering a profile allows applications to subscribe to *remote* services.
// These must implement the GattProfile1 interface defined above.
// Registering a service allows applications to publish a *local* GATT service,
// which then becomes available to remote devices. A GATT service is represented by
// a D-Bus object hierarchy where the root node corresponds to a service and the
// child nodes represent characteristics and descriptors that belong to that
// service. Each node must implement one of GattService1, GattCharacteristic1,
// or GattDescriptor1 interfaces described above, based on the attribute it
// represents. Each node must also implement the standard D-Bus Properties
// interface to expose their properties. These objects collectively represent a
// GATT service definition.
// To make service registration simple, BlueZ requires that all objects that belong
// to a GATT service be grouped under a D-Bus Object Manager that solely manages
// the objects of that service. Hence, the standard DBus.ObjectManager interface
// must be available on the root service path. An example application hierarchy
// containing two separate GATT services may look like this:
// -> /com/example
// |   - org.freedesktop.DBus.ObjectManager
// |
// -> /com/example/service0
// | |   - org.freedesktop.DBus.Properties
// | |   - org.bluez.GattService1
// | |
// | -> /com/example/service0/char0
// | |     - org.freedesktop.DBus.Properties
// | |     - org.bluez.GattCharacteristic1
// | |
// | -> /com/example/service0/char1
// |   |   - org.freedesktop.DBus.Properties
// |   |   - org.bluez.GattCharacteristic1
// |   |
// |   -> /com/example/service0/char1/desc0
// |       - org.freedesktop.DBus.Properties
// |       - org.bluez.GattDescriptor1
// |
// -> /com/example/service1
// |   - org.freedesktop.DBus.Properties
// |   - org.bluez.GattService1
// |
// -> /com/example/service1/char0
// - org.freedesktop.DBus.Properties
// - org.bluez.GattCharacteristic1
// When a service is registered, BlueZ will automatically obtain information about
// all objects using the service's Object Manager. Once a service has been
// registered, the objects of a service should not be removed. If BlueZ receives an
// InterfacesRemoved signal from a service's Object Manager, it will immediately
// unregister the service. Similarly, if the application disconnects from the bus,
// all of its registered services will be automatically unregistered.
// InterfacesAdded signals will be ignored.
// Examples:
// - Client
// test/example-gatt-client
// client/bluetoothctl
// - Server
// test/example-gatt-server
// tools/gatt-service
type GattManager1 struct {
	client     *bluez.Client
	Properties *GattManager1Properties
}

// GattManager1Properties contains the exposed properties of an interface
type GattManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *GattManager1Properties) Lock() {
	p.lock.Lock()
}

func (p *GattManager1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *GattManager1) Close() {
	a.client.Disconnect()
}


// ToMap convert a GattManager1Properties to map
func (a *GattManager1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an GattManager1Properties
func (a *GattManager1Properties) FromMap(props map[string]interface{}) (*GattManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattManager1Properties
func (a *GattManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattManager1Properties, error) {
	s := new(GattManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *GattManager1) GetProperties() (*GattManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *GattManager1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *GattManager1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//RegisterApplication Registers a local GATT services hierarchy as described
// above (GATT Server) and/or GATT profiles (GATT Client).
// The application object path together with the D-Bus
// system bus connection ID define the identification of
// the application registering a GATT based
// service or profile.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.AlreadyExists
func (a *GattManager1) RegisterApplication(application dbus.ObjectPath, options map[string]dbus.Variant) error {
	
	return a.client.Call("RegisterApplication", 0, application, options).Store()
	
}

//UnregisterApplication This unregisters the services that has been
// previously registered. The object path parameter
// must match the same value that has been used
// on registration.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.DoesNotExist
func (a *GattManager1) UnregisterApplication(application dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterApplication", 0, application).Store()
	
}

