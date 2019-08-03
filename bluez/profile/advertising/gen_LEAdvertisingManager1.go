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

package advertising



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
  "fmt"
)

var LEAdvertisingManager1Interface = "org.bluez.LEAdvertisingManager1"


// NewLEAdvertisingManager1 create a new instance of LEAdvertisingManager1
//
// Args:
// 	objectPath: /org/bluez/{hci0,hci1,...}
func NewLEAdvertisingManager1(objectPath string) (*LEAdvertisingManager1, error) {
	a := new(LEAdvertisingManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: LEAdvertisingManager1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(LEAdvertisingManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}

// NewLEAdvertisingManager1FromAdapterID create a new instance of LEAdvertisingManager1
// adapterID: ID of an adapter eg. hci0
func NewLEAdvertisingManager1FromAdapterID(adapterID string) (*LEAdvertisingManager1, error) {
	a := new(LEAdvertisingManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: LEAdvertisingManager1Interface,
			Path:  fmt.Sprintf("/org/bluez/%s", adapterID),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(LEAdvertisingManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// LEAdvertisingManager1 LE Advertising Manager hierarchy
// The Advertising Manager allows external applications to register Advertisement
// Data which should be broadcast to devices.  Advertisement Data elements must
// follow the API for LE Advertisement Data described above.
type LEAdvertisingManager1 struct {
	client     *bluez.Client
	Properties *LEAdvertisingManager1Properties
}

// LEAdvertisingManager1Properties contains the exposed properties of an interface
type LEAdvertisingManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// SupportedIncludes List of supported system includes.
  // Possible values: "tx-power"
  // "appearance"
  // "local-name"
	SupportedIncludes []string

	// ActiveInstances Number of active advertising instances.
	ActiveInstances byte

	// SupportedInstances Number of available advertising instances.
	SupportedInstances byte

}

func (p *LEAdvertisingManager1Properties) Lock() {
	p.lock.Lock()
}

func (p *LEAdvertisingManager1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *LEAdvertisingManager1) Close() {
	a.client.Disconnect()
}


// ToMap convert a LEAdvertisingManager1Properties to map
func (a *LEAdvertisingManager1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an LEAdvertisingManager1Properties
func (a *LEAdvertisingManager1Properties) FromMap(props map[string]interface{}) (*LEAdvertisingManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an LEAdvertisingManager1Properties
func (a *LEAdvertisingManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*LEAdvertisingManager1Properties, error) {
	s := new(LEAdvertisingManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *LEAdvertisingManager1) GetProperties() (*LEAdvertisingManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *LEAdvertisingManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *LEAdvertisingManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *LEAdvertisingManager1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *LEAdvertisingManager1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//RegisterAdvertisement Registers an advertisement object to be sent over the LE
// Advertising channel.  The service must be exported
// under interface LEAdvertisement1.
// InvalidArguments error indicates that the object has
// invalid or conflicting properties.
// InvalidLength error indicates that the data
// provided generates a data packet which is too long.
// The properties of this object are parsed when it is
// registered, and any changes are ignored.
// If the same object is registered twice it will result in
// an AlreadyExists error.
// If the maximum number of advertisement instances is
// reached it will result in NotPermitted error.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.AlreadyExists
// org.bluez.Error.InvalidLength
func (a *LEAdvertisingManager1) RegisterAdvertisement(advertisement dbus.ObjectPath, options map[string]interface{}) error {
	
	return a.client.Call("RegisterAdvertisement", 0, advertisement, options).Store()
	
}

//UnregisterAdvertisement This unregisters an advertisement that has been
// previously registered.  The object path parameter must
// match the same value that has been used on registration.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.DoesNotExist
func (a *LEAdvertisingManager1) UnregisterAdvertisement(advertisement dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterAdvertisement", 0, advertisement).Store()
	
}

