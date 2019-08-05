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

var HealthChannel1Interface = "org.bluez.HealthChannel1"


// NewHealthChannel1 create a new instance of HealthChannel1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/chanZZZ
func NewHealthChannel1(objectPath dbus.ObjectPath) (*HealthChannel1, error) {
	a := new(HealthChannel1)
	a.propertiesSignal = make(chan *dbus.Signal)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: HealthChannel1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(HealthChannel1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// HealthChannel1 HealthChannel hierarchy

type HealthChannel1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	Properties 				*HealthChannel1Properties
}

// HealthChannel1Properties contains the exposed properties of an interface
type HealthChannel1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Application Identifies the HealthApplication to which this channel
  // is related to (which indirectly defines its role and
  // data type).
	Application dbus.ObjectPath

	// Type The quality of service of the data channel. ("reliable"
  // or "streaming")
	Type string

	// Device Identifies the Remote Device that is connected with.
  // Maps with a HealthDevice object.
	Device dbus.ObjectPath

}

func (p *HealthChannel1Properties) Lock() {
	p.lock.Lock()
}

func (p *HealthChannel1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *HealthChannel1) Close() {
	
	a.unregisterSignal()
	
	a.client.Disconnect()
}

// Path return HealthChannel1 object path
func (a *HealthChannel1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return HealthChannel1 interface
func (a *HealthChannel1) Interface() string {
	return a.client.Config.Iface
}


// ToMap convert a HealthChannel1Properties to map
func (a *HealthChannel1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an HealthChannel1Properties
func (a *HealthChannel1Properties) FromMap(props map[string]interface{}) (*HealthChannel1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an HealthChannel1Properties
func (a *HealthChannel1Properties) FromDBusMap(props map[string]dbus.Variant) (*HealthChannel1Properties, error) {
	s := new(HealthChannel1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *HealthChannel1) GetProperties() (*HealthChannel1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *HealthChannel1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *HealthChannel1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *HealthChannel1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *HealthChannel1) unregisterSignal() {
	if a.propertiesSignal == nil {
		a.propertiesSignal <- nil
	}
}

// WatchProperties updates on property changes
func (a *HealthChannel1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *HealthChannel1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}





//Acquire Returns the file descriptor for this data channel. If
// the data channel is not connected it will also
// reconnect.
// Possible Errors: org.bluez.Error.NotConnected
// org.bluez.Error.NotAllowed
func (a *HealthChannel1) Acquire() (dbus.UnixFD, error) {
	
	var val0 dbus.UnixFD
	err := a.client.Call("Acquire", 0, ).Store(&val0)
	return val0, err	
}

//Release 
func (a *HealthChannel1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

//close Possible Errors: org.bluez.Error.NotAcquired
// org.bluez.Error.NotAllowed
func (a *HealthChannel1) close() error {
	
	return a.client.Call("close", 0, ).Store()
	
}

