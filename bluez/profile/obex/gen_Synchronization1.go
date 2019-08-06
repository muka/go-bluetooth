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

package obex



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Synchronization1Interface = "org.bluez.obex.Synchronization1"


// NewSynchronization1 create a new instance of Synchronization1
//
// Args:
// 	objectPath: [Session object path]
func NewSynchronization1(objectPath dbus.ObjectPath) (*Synchronization1, error) {
	a := new(Synchronization1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: Synchronization1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Synchronization1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Synchronization1 Synchronization hierarchy

type Synchronization1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Synchronization1Properties
}

// Synchronization1Properties contains the exposed properties of an interface
type Synchronization1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *Synchronization1Properties) Lock() {
	p.lock.Lock()
}

func (p *Synchronization1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *Synchronization1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return Synchronization1 object path
func (a *Synchronization1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return Synchronization1 interface
func (a *Synchronization1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Synchronization1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a Synchronization1Properties to map
func (a *Synchronization1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Synchronization1Properties
func (a *Synchronization1Properties) FromMap(props map[string]interface{}) (*Synchronization1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Synchronization1Properties
func (a *Synchronization1Properties) FromDBusMap(props map[string]dbus.Variant) (*Synchronization1Properties, error) {
	s := new(Synchronization1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Synchronization1) GetProperties() (*Synchronization1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Synchronization1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Synchronization1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Synchronization1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Synchronization1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *Synchronization1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *Synchronization1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




//SetLocation Set the phonebook object store location for other
// operations. Should be called before all the other
// operations.
// location: Where the phonebook is stored, possible
// values:
// "int" ( "internal" which is default )
// "sim1"
// "sim2"
// ......
// Possible errors: org.bluez.obex.Error.InvalidArguments
func (a *Synchronization1) SetLocation(location string) error {
	
	return a.client.Call("SetLocation", 0, location).Store()
	
}

//GetPhonebook Retrieve an entire Phonebook Object store from remote
// device, and stores it in a local file.
// If an empty target file is given, a name will be
// automatically calculated for the temporary file.
// The returned path represents the newly created transfer,
// which should be used to find out if the content has been
// successfully transferred or if the operation fails.
// The properties of this transfer are also returned along
// with the object path, to avoid a call to GetProperties.
// Possible errors: org.bluez.obex.Error.InvalidArguments
// org.bluez.obex.Error.Failed
func (a *Synchronization1) GetPhonebook(targetfile string) (dbus.ObjectPath, map[string]interface{}, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]interface{}
	err := a.client.Call("GetPhonebook", 0, targetfile).Store(&val0, &val1)
	return val0, val1, err	
}

//PutPhonebook Send an entire Phonebook Object store to remote device.
// The returned path represents the newly created transfer,
// which should be used to find out if the content has been
// successfully transferred or if the operation fails.
// The properties of this transfer are also returned along
// with the object path, to avoid a call to GetProperties.
// Possible errors: org.bluez.obex.Error.InvalidArguments
// org.bluez.obex.Error.Failed
func (a *Synchronization1) PutPhonebook(sourcefile string) (dbus.ObjectPath, map[string]interface{}, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]interface{}
	err := a.client.Call("PutPhonebook", 0, sourcefile).Store(&val0, &val1)
	return val0, val1, err	
}

