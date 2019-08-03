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

var GattDescriptor1Interface = "org.bluez.GattDescriptor1"


// NewGattDescriptor1 create a new instance of GattDescriptor1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/serviceXX/charYYYY/descriptorZZZ
func NewGattDescriptor1(objectPath string) (*GattDescriptor1, error) {
	a := new(GattDescriptor1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattDescriptor1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattDescriptor1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// GattDescriptor1 Characteristic Descriptors hierarchy
// Local or remote GATT characteristic descriptors hierarchy.
type GattDescriptor1 struct {
	client     *bluez.Client
	Properties *GattDescriptor1Properties
}

// GattDescriptor1Properties contains the exposed properties of an interface
type GattDescriptor1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Flags Defines how the descriptor value can be used.
  // Possible values:
  // "read"
  // "write"
  // "encrypt-read"
  // "encrypt-write"
  // "encrypt-authenticated-read"
  // "encrypt-authenticated-write"
  // "secure-read" (Server Only)
  // "secure-write" (Server Only)
  // "authorize"
	Flags []string

	// UUID 128-bit descriptor UUID.
	UUID string

	// Characteristic Object path of the GATT characteristic the descriptor
  // belongs to.
	Characteristic dbus.ObjectPath

	// Value The cached value of the descriptor. This property
  // gets updated only after a successful read request, upon
  // which a PropertiesChanged signal will be emitted.
	Value []byte `dbus:"emit"`

}

func (p *GattDescriptor1Properties) Lock() {
	p.lock.Lock()
}

func (p *GattDescriptor1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *GattDescriptor1) Close() {
	a.client.Disconnect()
}


// ToMap convert a GattDescriptor1Properties to map
func (a *GattDescriptor1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an GattDescriptor1Properties
func (a *GattDescriptor1Properties) FromMap(props map[string]interface{}) (*GattDescriptor1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattDescriptor1Properties
func (a *GattDescriptor1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattDescriptor1Properties, error) {
	s := new(GattDescriptor1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *GattDescriptor1) GetProperties() (*GattDescriptor1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattDescriptor1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattDescriptor1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *GattDescriptor1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *GattDescriptor1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//ReadValue Issues a request to read the value of the
// characteristic and returns the value if the
// operation was successful.
// Possible options: "offset": Start offset
// "device": Device path (Server only)
// "link": Link type (Server only)
// Possible Errors: org.bluez.Error.Failed
// org.bluez.Error.InProgress
// org.bluez.Error.NotPermitted
// org.bluez.Error.NotAuthorized
// org.bluez.Error.NotSupported
func (a *GattDescriptor1) ReadValue(flags map[string]interface{}) ([]byte, error) {
	
	var val0 []byte
	err := a.client.Call("ReadValue", 0, flags).Store(&val0)
	return val0, err	
}

//WriteValue Issues a request to write the value of the
// characteristic.
// Possible options: "offset": Start offset
// "device": Device path (Server only)
// "link": Link type (Server only)
// "prepare-authorize": boolean Is prepare
// authorization
// request
// Possible Errors: org.bluez.Error.Failed
// org.bluez.Error.InProgress
// org.bluez.Error.NotPermitted
// org.bluez.Error.InvalidValueLength
// org.bluez.Error.NotAuthorized
// org.bluez.Error.NotSupported
func (a *GattDescriptor1) WriteValue(value []byte, flags map[string]interface{}) error {
	
	return a.client.Call("WriteValue", 0, value, flags).Store()
	
}

