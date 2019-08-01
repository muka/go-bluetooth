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

package media



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var MediaEndpoint1Interface = "org.bluez.MediaEndpoint1"


// NewMediaEndpoint1 create a new instance of MediaEndpoint1
//
// Args:
// 	servicePath: unique name
// 	objectPath: freely definable
func NewMediaEndpoint1(servicePath string, objectPath string) (*MediaEndpoint1, error) {
	a := new(MediaEndpoint1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: MediaEndpoint1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaEndpoint1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// MediaEndpoint1 MediaEndpoint1 hierarchy

type MediaEndpoint1 struct {
	client     *bluez.Client
	Properties *MediaEndpoint1Properties
}

// MediaEndpoint1Properties contains the exposed properties of an interface
type MediaEndpoint1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *MediaEndpoint1Properties) Lock() {
	p.lock.Lock()
}

func (p *MediaEndpoint1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *MediaEndpoint1) Close() {
	a.client.Disconnect()
}


// ToMap convert a MediaEndpoint1Properties to map
func (a *MediaEndpoint1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MediaEndpoint1Properties
func (a *MediaEndpoint1Properties) FromMap(props map[string]interface{}) (*MediaEndpoint1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MediaEndpoint1Properties
func (a *MediaEndpoint1Properties) FromDBusMap(props map[string]dbus.Variant) (*MediaEndpoint1Properties, error) {
	s := new(MediaEndpoint1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MediaEndpoint1) GetProperties() (*MediaEndpoint1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MediaEndpoint1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MediaEndpoint1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *MediaEndpoint1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *MediaEndpoint1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//SetConfiguration Set configuration for the transport.
func (a *MediaEndpoint1) SetConfiguration(transport dbus.ObjectPath, properties map[string]dbus.Variant) error {
	
	return a.client.Call("SetConfiguration", 0, transport, properties).Store()
	
}

//SelectConfiguration Select preferable configuration from the supported
// capabilities.
// Returns a configuration which can be used to setup
// a transport.
// Note: There is no need to cache the selected
// configuration since on success the configuration is
// send back as parameter of SetConfiguration.
func (a *MediaEndpoint1) SelectConfiguration(capabilities []byte) ([]byte, error) {
	
	var val0 []byte
	err := a.client.Call("SelectConfiguration", 0, capabilities).Store(&val0)
	return val0, err	
}

//ClearConfiguration Clear transport configuration.
func (a *MediaEndpoint1) ClearConfiguration(transport dbus.ObjectPath) error {
	
	return a.client.Call("ClearConfiguration", 0, transport).Store()
	
}

//Release This method gets called when the service daemon
// unregisters the endpoint. An endpoint can use it to do
// cleanup tasks. There is no need to unregister the
// endpoint, because when this method gets called it has
// already been unregistered.
func (a *MediaEndpoint1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

