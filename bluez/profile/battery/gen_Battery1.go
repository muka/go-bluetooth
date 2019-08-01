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

package battery



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Battery1Interface = "org.bluez.Battery1"


// NewBattery1 create a new instance of Battery1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewBattery1(objectPath string) (*Battery1, error) {
	a := new(Battery1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Battery1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Battery1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Battery1 Battery hierarchy

type Battery1 struct {
	client     *bluez.Client
	Properties *Battery1Properties
}

// Battery1Properties contains the exposed properties of an interface
type Battery1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Percentage The percentage of battery left as an unsigned 8-bit integer.
	Percentage byte

}

func (p *Battery1Properties) Lock() {
	p.lock.Lock()
}

func (p *Battery1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Battery1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Battery1Properties to map
func (a *Battery1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Battery1Properties
func (a *Battery1Properties) FromMap(props map[string]interface{}) (*Battery1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Battery1Properties
func (a *Battery1Properties) FromDBusMap(props map[string]dbus.Variant) (*Battery1Properties, error) {
	s := new(Battery1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Battery1) GetProperties() (*Battery1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Battery1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Battery1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Battery1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Battery1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



