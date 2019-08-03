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

package profile



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var ProfileManager1Interface = "org.bluez.ProfileManager1"


// NewProfileManager1 create a new instance of ProfileManager1
//
// Args:

func NewProfileManager1() (*ProfileManager1, error) {
	a := new(ProfileManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: ProfileManager1Interface,
			Path:  "/org/bluez",
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(ProfileManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// ProfileManager1 Profile Manager hierarchy

type ProfileManager1 struct {
	client     *bluez.Client
	Properties *ProfileManager1Properties
}

// ProfileManager1Properties contains the exposed properties of an interface
type ProfileManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *ProfileManager1Properties) Lock() {
	p.lock.Lock()
}

func (p *ProfileManager1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *ProfileManager1) Close() {
	a.client.Disconnect()
}


// ToMap convert a ProfileManager1Properties to map
func (a *ProfileManager1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an ProfileManager1Properties
func (a *ProfileManager1Properties) FromMap(props map[string]interface{}) (*ProfileManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an ProfileManager1Properties
func (a *ProfileManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*ProfileManager1Properties, error) {
	s := new(ProfileManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *ProfileManager1) GetProperties() (*ProfileManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *ProfileManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *ProfileManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *ProfileManager1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *ProfileManager1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//RegisterProfile This registers a profile implementation.
// If an application disconnects from the bus all
// its registered profiles will be removed.
// HFP HS UUID: 0000111e-0000-1000-8000-00805f9b34fb
// Default RFCOMM channel is 6. And this requires
// authentication.
// Available options:
// string Name
// Human readable name for the profile
// string Service
func (a *ProfileManager1) RegisterProfile(profile dbus.ObjectPath, uuid string, options map[string]interface{}) error {
	
	return a.client.Call("RegisterProfile", 0, profile, uuid, options).Store()
	
}

//UnregisterProfile This unregisters the profile that has been previously
// registered. The object path parameter must match the
// same value that has been used on registration.
// Possible errors: org.bluez.Error.DoesNotExist
func (a *ProfileManager1) UnregisterProfile(profile dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterProfile", 0, profile).Store()
	
}

