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

package input



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Input1Interface = "org.bluez.Input1"


// NewInput1 create a new instance of Input1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewInput1(objectPath string) (*Input1, error) {
	a := new(Input1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Input1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Input1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Input1 Input hierarchy

type Input1 struct {
	client     *bluez.Client
	Properties *Input1Properties
}

// Input1Properties contains the exposed properties of an interface
type Input1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// ReconnectMode Determines the Connectability mode of the HID device as
  // defined by the HID Profile specification, Section 5.4.2.
  // This mode is based in the two properties
  // HIDReconnectInitiate (see Section 5.3.4.6) and
  // HIDNormallyConnectable (see Section 5.3.4.14) which
  // define the following four possible values:
  // "none"		Device and host are not required to
  // automatically restore the connection.
  // "host"		Bluetooth HID host restores connection.
  // "device"	Bluetooth HID device restores
  // connection.
  // "any"		Bluetooth HID device shall attempt to
  // restore the lost connection, but
  // Bluetooth HID Host may also restore the
  // connection.
	ReconnectMode string

}

func (p *Input1Properties) Lock() {
	p.lock.Lock()
}

func (p *Input1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Input1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Input1Properties to map
func (a *Input1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Input1Properties
func (a *Input1Properties) FromMap(props map[string]interface{}) (*Input1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Input1Properties
func (a *Input1Properties) FromDBusMap(props map[string]dbus.Variant) (*Input1Properties, error) {
	s := new(Input1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Input1) GetProperties() (*Input1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Input1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Input1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Input1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Input1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



