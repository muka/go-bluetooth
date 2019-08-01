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
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Synchronization1Interface = "org.bluez.obex.Synchronization1"


// NewSynchronization1 create a new instance of Synchronization1
//
// Args:
// 	objectPath: [Session object path]
func NewSynchronization1(objectPath string) (*Synchronization1, error) {
	a := new(Synchronization1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: Synchronization1Interface,
			Path:  objectPath,
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
	client     *bluez.Client
	Properties *Synchronization1Properties
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
	a.client.Disconnect()
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

// Register for changes signalling
func (a *Synchronization1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Synchronization1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
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
func (a *Synchronization1) GetPhonebook(targetfile string) (dbus.ObjectPath, map[string]dbus.Variant, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]dbus.Variant
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
func (a *Synchronization1) PutPhonebook(sourcefile string) (dbus.ObjectPath, map[string]dbus.Variant, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]dbus.Variant
	err := a.client.Call("PutPhonebook", 0, sourcefile).Store(&val0, &val1)
	return val0, val1, err	
}

