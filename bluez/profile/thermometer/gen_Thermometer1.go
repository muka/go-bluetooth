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
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Thermometer1Interface = "org.bluez.Thermometer1"


// NewThermometer1 create a new instance of Thermometer1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewThermometer1(objectPath dbus.ObjectPath) (*Thermometer1, error) {
	a := new(Thermometer1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Thermometer1Interface,
			Path:  dbus.ObjectPath(objectPath),
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
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Thermometer1Properties
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


// SetIntermediate set Intermediate value
func (a *Thermometer1) SetIntermediate(v bool) error {
	return a.SetProperty("Intermediate", v)
}

// GetIntermediate get Intermediate value
func (a *Thermometer1) GetIntermediate() (bool, error) {
	v, err := a.GetProperty("Intermediate")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetInterval set Interval value
func (a *Thermometer1) SetInterval(v uint16) error {
	return a.SetProperty("Interval", v)
}

// GetInterval get Interval value
func (a *Thermometer1) GetInterval() (uint16, error) {
	v, err := a.GetProperty("Interval")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}

// SetMaximum set Maximum value
func (a *Thermometer1) SetMaximum(v uint16) error {
	return a.SetProperty("Maximum", v)
}

// GetMaximum get Maximum value
func (a *Thermometer1) GetMaximum() (uint16, error) {
	v, err := a.GetProperty("Maximum")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}

// SetMinimum set Minimum value
func (a *Thermometer1) SetMinimum(v uint16) error {
	return a.SetProperty("Minimum", v)
}

// GetMinimum get Minimum value
func (a *Thermometer1) GetMinimum() (uint16, error) {
	v, err := a.GetProperty("Minimum")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}


// Close the connection
func (a *Thermometer1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return Thermometer1 object path
func (a *Thermometer1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return Thermometer1 interface
func (a *Thermometer1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Thermometer1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Thermometer1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Thermometer1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *Thermometer1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *Thermometer1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




