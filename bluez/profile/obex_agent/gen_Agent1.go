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

package obex_agent



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Agent1Interface = "org.bluez.obex.Agent1"


// NewAgent1 create a new instance of Agent1
//
// Args:
// 	servicePath: unique name
// 	objectPath: freely definable
func NewAgent1(servicePath string, objectPath string) (*Agent1, error) {
	a := new(Agent1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: Agent1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Agent1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Agent1 Agent hierarchy

type Agent1 struct {
	client     *bluez.Client
	Properties *Agent1Properties
}

// Agent1Properties contains the exposed properties of an interface
type Agent1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *Agent1Properties) Lock() {
	p.lock.Lock()
}

func (p *Agent1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Agent1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Agent1Properties to map
func (a *Agent1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Agent1Properties
func (a *Agent1Properties) FromMap(props map[string]interface{}) (*Agent1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Agent1Properties
func (a *Agent1Properties) FromDBusMap(props map[string]dbus.Variant) (*Agent1Properties, error) {
	s := new(Agent1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Agent1) GetProperties() (*Agent1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Agent1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Agent1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Agent1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Agent1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Release This method gets called when the service daemon
// unregisters the agent. An agent can use it to do
// cleanup tasks. There is no need to unregister the
// agent, because when this method gets called it has
// already been unregistered.
func (a *Agent1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

//AuthorizePush This method gets called when the service daemon
// needs to accept/reject a Bluetooth object push request.
// Returns the full path (including the filename) where
// the object shall be stored. The tranfer object will
// contain a Filename property that contains the default
// location and name that can be returned.
// Possible errors: org.bluez.obex.Error.Rejected
// org.bluez.obex.Error.Canceled
func (a *Agent1) AuthorizePush(transfer dbus.ObjectPath) (string, error) {
	
	var val0 string
	err := a.client.Call("AuthorizePush", 0, transfer).Store(&val0)
	return val0, err	
}

//Cancel This method gets called to indicate that the agent
// request failed before a reply was returned. It cancels
// the previous request.
func (a *Agent1) Cancel() error {
	
	return a.client.Call("Cancel", 0, ).Store()
	
}

