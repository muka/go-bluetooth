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

package health



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var HealthManager1Interface = "org.bluez.HealthManager1"


// NewHealthManager1 create a new instance of HealthManager1
//
// Args:

func NewHealthManager1() (*HealthManager1, error) {
	a := new(HealthManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: HealthManager1Interface,
			Path:  "/org/bluez/",
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(HealthManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// HealthManager1 HealthManager hierarchy

type HealthManager1 struct {
	client     *bluez.Client
	Properties *HealthManager1Properties
}

// HealthManager1Properties contains the exposed properties of an interface
type HealthManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *HealthManager1Properties) Lock() {
	p.lock.Lock()
}

func (p *HealthManager1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *HealthManager1) Close() {
	a.client.Disconnect()
}


// ToMap convert a HealthManager1Properties to map
func (a *HealthManager1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an HealthManager1Properties
func (a *HealthManager1Properties) FromMap(props map[string]interface{}) (*HealthManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an HealthManager1Properties
func (a *HealthManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*HealthManager1Properties, error) {
	s := new(HealthManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *HealthManager1) GetProperties() (*HealthManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *HealthManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *HealthManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *HealthManager1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *HealthManager1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//CreateApplication Returns the path of the new registered application.
// Application will be closed by the call or implicitly
// when the programs leaves the bus.
// config:
// uint16 DataType:
// Mandatory
// string Role:
// Mandatory. Possible values: "source",
// "sink"
// string Description:
// Optional
// ChannelType:
// Optional, just for sources. Possible
// values: "reliable", "streaming"
// Possible Errors: org.bluez.Error.InvalidArguments
func (a *HealthManager1) CreateApplication(config map[string]interface{}) (dbus.ObjectPath, error) {
	
	var val0 dbus.ObjectPath
	err := a.client.Call("CreateApplication", 0, config).Store(&val0)
	return val0, err	
}

//DestroyApplication Closes the HDP application identified by the object
// path. Also application will be closed if the process
// that started it leaves the bus. Only the creator of the
// application will be able to destroy it.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.NotFound
// org.bluez.Error.NotAllowed
func (a *HealthManager1) DestroyApplication(application dbus.ObjectPath) error {
	
	return a.client.Call("DestroyApplication", 0, application).Store()
	
}

