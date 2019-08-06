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

// GATT remote and local service representation. Object path for local services
// is freely definable.
// External applications implementing local services must register the services
// using GattManager1 registration method and must implement the methods and
// properties defined in GattService1 interface.
package gatt



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var GattService1Interface = "org.bluez.GattService1"


// NewGattService1 create a new instance of GattService1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/serviceXX
func NewGattService1(objectPath dbus.ObjectPath) (*GattService1, error) {
	a := new(GattService1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattService1Interface,
			Path:  dbus.ObjectPath(objectPath),
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
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*GattService1Properties
}

// GattService1Properties contains the exposed properties of an interface
type GattService1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

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

	// Characteristics 
	Characteristics []dbus.ObjectPath `dbus:"emit"`

	// IsService 
	IsService bool `dbus:"ignore"`

}

func (p *GattService1Properties) Lock() {
	p.lock.Lock()
}

func (p *GattService1Properties) Unlock() {
	p.lock.Unlock()
}


// SetUUID set UUID value
func (a *GattService1) SetUUID(v string) error {
	return a.SetProperty("UUID", v)
}

// GetUUID get UUID value
func (a *GattService1) GetUUID() (string, error) {
	v, err := a.GetProperty("UUID")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetPrimary set Primary value
func (a *GattService1) SetPrimary(v bool) error {
	return a.SetProperty("Primary", v)
}

// GetPrimary get Primary value
func (a *GattService1) GetPrimary() (bool, error) {
	v, err := a.GetProperty("Primary")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetDevice set Device value
func (a *GattService1) SetDevice(v dbus.ObjectPath) error {
	return a.SetProperty("Device", v)
}

// GetDevice get Device value
func (a *GattService1) GetDevice() (dbus.ObjectPath, error) {
	v, err := a.GetProperty("Device")
	if err != nil {
		return dbus.ObjectPath(""), err
	}
	return v.Value().(dbus.ObjectPath), nil
}

// SetIncludes set Includes value
func (a *GattService1) SetIncludes(v []dbus.ObjectPath) error {
	return a.SetProperty("Includes", v)
}

// GetIncludes get Includes value
func (a *GattService1) GetIncludes() ([]dbus.ObjectPath, error) {
	v, err := a.GetProperty("Includes")
	if err != nil {
		return []dbus.ObjectPath{}, err
	}
	return v.Value().([]dbus.ObjectPath), nil
}

// SetCharacteristics set Characteristics value
func (a *GattService1) SetCharacteristics(v []dbus.ObjectPath) error {
	return a.SetProperty("Characteristics", v)
}

// GetCharacteristics get Characteristics value
func (a *GattService1) GetCharacteristics() ([]dbus.ObjectPath, error) {
	v, err := a.GetProperty("Characteristics")
	if err != nil {
		return []dbus.ObjectPath{}, err
	}
	return v.Value().([]dbus.ObjectPath), nil
}

// SetIsService set IsService value
func (a *GattService1) SetIsService(v bool) error {
	return a.SetProperty("IsService", v)
}

// GetIsService get IsService value
func (a *GattService1) GetIsService() (bool, error) {
	v, err := a.GetProperty("IsService")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}


// Close the connection
func (a *GattService1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return GattService1 object path
func (a *GattService1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return GattService1 interface
func (a *GattService1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *GattService1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
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

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *GattService1) GetPropertiesSignal() (chan *dbus.Signal, error) {

	if a.propertiesSignal == nil {
		s, err := a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
		if err != nil {
			return nil, err
		}
		a.propertiesSignal = s
	}

	return a.propertiesSignal, nil
}

// Unregister for changes signalling
func (a *GattService1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *GattService1) WatchProperties() (chan *bluez.PropertyChanged, error) {

	channel, err := a.client.Register(a.Path(), a.Interface())
	if err != nil {
		return nil, err
	}

	ch := make(chan *bluez.PropertyChanged)

	go (func() {
		for {

			if channel == nil {
				break
			}

			sig := <-channel

			if sig == nil {
				return
			}

			if sig.Name != bluez.PropertiesChanged {
				continue
			}
			if sig.Path != a.Path() {
				continue
			}

			iface := sig.Body[0].(string)
			changes := sig.Body[1].(map[string]dbus.Variant)

			for field, val := range changes {

				// updates [*]Properties struct when a property change
				s := reflect.ValueOf(a.Properties).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						a.Properties.Lock()
						f.Set(x)
						a.Properties.Unlock()
					}
				}

				propChanged := &bluez.PropertyChanged{
					Interface: iface,
					Name:      field,
					Value:     val.Value(),
				}
				ch <- propChanged
			}

		}
	})()

	return ch, nil
}

func (a *GattService1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




