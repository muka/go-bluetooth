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
  "reflect"
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
			Path:  dbus.ObjectPath("/org/bluez/"),
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
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager	
	Properties 				*HealthManager1Properties
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
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return HealthManager1 object path
func (a *HealthManager1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return HealthManager1 interface
func (a *HealthManager1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *HealthManager1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *HealthManager1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *HealthManager1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *HealthManager1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

				// updates [*]Properties struct
				props := a.Properties

				s := reflect.ValueOf(props).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						props.Lock()
						f.Set(x)
						props.Unlock()
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

func (a *HealthManager1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
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

