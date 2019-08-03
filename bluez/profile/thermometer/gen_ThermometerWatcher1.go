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

var ThermometerWatcher1Interface = "org.bluez.ThermometerWatcher1"


// NewThermometerWatcher1 create a new instance of ThermometerWatcher1
//
// Args:
// 	servicePath: unique name
// 	objectPath: freely definable
func NewThermometerWatcher1(servicePath string, objectPath string) (*ThermometerWatcher1, error) {
	a := new(ThermometerWatcher1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: ThermometerWatcher1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(ThermometerWatcher1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// ThermometerWatcher1 Health Thermometer Watcher hierarchy

type ThermometerWatcher1 struct {
	client     *bluez.Client
	Properties *ThermometerWatcher1Properties
}

// ThermometerWatcher1Properties contains the exposed properties of an interface
type ThermometerWatcher1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *ThermometerWatcher1Properties) Lock() {
	p.lock.Lock()
}

func (p *ThermometerWatcher1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *ThermometerWatcher1) Close() {
	a.client.Disconnect()
}


// ToMap convert a ThermometerWatcher1Properties to map
func (a *ThermometerWatcher1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an ThermometerWatcher1Properties
func (a *ThermometerWatcher1Properties) FromMap(props map[string]interface{}) (*ThermometerWatcher1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an ThermometerWatcher1Properties
func (a *ThermometerWatcher1Properties) FromDBusMap(props map[string]dbus.Variant) (*ThermometerWatcher1Properties, error) {
	s := new(ThermometerWatcher1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *ThermometerWatcher1) GetProperties() (*ThermometerWatcher1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *ThermometerWatcher1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *ThermometerWatcher1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *ThermometerWatcher1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *ThermometerWatcher1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//MeasurementReceived This callback gets called when a measurement has been
// scanned in the thermometer.
// Measurement:
// int16 Exponent:
// int32 Mantissa:
// Exponent and Mantissa values as
// extracted from float value defined by
// IEEE-11073-20601.
func (a *ThermometerWatcher1) MeasurementReceived(measurement map[string]interface{}) error {
	
	return a.client.Call("MeasurementReceived", 0, measurement).Store()
	
}

