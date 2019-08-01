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

var Profile1Interface = "org.bluez.Profile1"


// NewProfile1 create a new instance of Profile1
//
// Args:
// 	servicePath: unique name
// 	objectPath: freely definable
func NewProfile1(servicePath string, objectPath string) (*Profile1, error) {
	a := new(Profile1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: Profile1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Profile1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Profile1 Profile hierarchy

type Profile1 struct {
	client     *bluez.Client
	Properties *Profile1Properties
}

// Profile1Properties contains the exposed properties of an interface
type Profile1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *Profile1Properties) Lock() {
	p.lock.Lock()
}

func (p *Profile1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Profile1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Profile1Properties to map
func (a *Profile1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Profile1Properties
func (a *Profile1Properties) FromMap(props map[string]interface{}) (*Profile1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Profile1Properties
func (a *Profile1Properties) FromDBusMap(props map[string]dbus.Variant) (*Profile1Properties, error) {
	s := new(Profile1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Profile1) GetProperties() (*Profile1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Profile1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Profile1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Profile1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Profile1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Release This method gets called when the service daemon
// unregisters the profile. A profile can use it to do
// cleanup tasks. There is no need to unregister the
// profile, because when this method gets called it has
// already been unregistered.
func (a *Profile1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

//NewConnection This method gets called when a new service level
// connection has been made and authorized.
// Common fd_properties:
// uint16 Version		Profile version (optional)
// uint16 Features		Profile features (optional)
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Profile1) NewConnection(device dbus.ObjectPath, fd int32, fd_properties map[string]dbus.Variant) error {
	
	return a.client.Call("NewConnection", 0, device, fd, fd_properties).Store()
	
}

//RequestDisconnection This method gets called when a profile gets
// disconnected.
// The file descriptor is no longer owned by the service
// daemon and the profile implementation needs to take
// care of cleaning up all connections.
// If multiple file descriptors are indicated via
// NewConnection, it is expected that all of them
// are disconnected before returning from this
// method call.
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Profile1) RequestDisconnection(device dbus.ObjectPath) error {
	
	return a.client.Call("RequestDisconnection", 0, device).Store()
	
}

