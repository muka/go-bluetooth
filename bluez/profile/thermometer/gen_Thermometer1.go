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

package thermometer



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Thermometer1Interface = "org.bluez.Thermometer1"


// NewThermometer1 create a new instance of Thermometer1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewThermometer1(objectPath string) (*Thermometer1, error) {
	a := new(Thermometer1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Thermometer1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Thermometer1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Thermometer1 Health Thermometer Profile hierarchy

type Thermometer1 struct {
	client     *bluez.Client
	Properties *Thermometer1Properties
}

// Thermometer1Properties contains the exposed properties of an interface
type Thermometer1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Intermediate True if the thermometer supports intermediate
  // measurement notifications.
	Intermediate bool

	// Interval (optional) The Measurement Interval defines the time (in
  // seconds) between measurements. This interval is
  // not related to the intermediate measurements and
  // must be defined into a valid range. Setting it
  // to zero means that no periodic measurements will
  // be taken.
	Interval uint16

	// Maximum (optional) Defines the maximum value allowed for the interval
  // between periodic measurements.
	Maximum uint16

	// Minimum (optional) Defines the minimum value allowed for the interval
  // between periodic measurements.
	Minimum uint16

}

func (p *Thermometer1Properties) Lock() {
	p.lock.Lock()
}

func (p *Thermometer1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Thermometer1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Thermometer1Properties to map
func (a *Thermometer1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Thermometer1Properties
func (a *Thermometer1Properties) FromMap(props map[string]interface{}) (*Thermometer1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Thermometer1Properties
func (a *Thermometer1Properties) FromDBusMap(props map[string]dbus.Variant) (*Thermometer1Properties, error) {
	s := new(Thermometer1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Thermometer1) GetProperties() (*Thermometer1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Thermometer1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Thermometer1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Thermometer1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Thermometer1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



