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

var AgentManager1Interface = "org.bluez.obex.AgentManager1"


// NewAgentManager1 create a new instance of AgentManager1
//
// Args:

func NewAgentManager1() (*AgentManager1, error) {
	a := new(AgentManager1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: AgentManager1Interface,
			Path:  "/org/bluez/obex",
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(AgentManager1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// AgentManager1 Agent Manager hierarchy

type AgentManager1 struct {
	client     *bluez.Client
	Properties *AgentManager1Properties
}

// AgentManager1Properties contains the exposed properties of an interface
type AgentManager1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *AgentManager1Properties) Lock() {
	p.lock.Lock()
}

func (p *AgentManager1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *AgentManager1) Close() {
	a.client.Disconnect()
}


// ToMap convert a AgentManager1Properties to map
func (a *AgentManager1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an AgentManager1Properties
func (a *AgentManager1Properties) FromMap(props map[string]interface{}) (*AgentManager1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an AgentManager1Properties
func (a *AgentManager1Properties) FromDBusMap(props map[string]dbus.Variant) (*AgentManager1Properties, error) {
	s := new(AgentManager1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *AgentManager1) GetProperties() (*AgentManager1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *AgentManager1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *AgentManager1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *AgentManager1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *AgentManager1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//RegisterAgent Register an agent to request authorization of
// the user to accept/reject objects. Object push
// service needs to authorize each received object.
// Possible errors: org.bluez.obex.Error.AlreadyExists
func (a *AgentManager1) RegisterAgent(agent dbus.ObjectPath) error {
	
	return a.client.Call("RegisterAgent", 0, agent).Store()
	
}

//UnregisterAgent This unregisters the agent that has been previously
// registered. The object path parameter must match the
// same value that has been used on registration.
// Possible errors: org.bluez.obex.Error.DoesNotExist
func (a *AgentManager1) UnregisterAgent(agent dbus.ObjectPath) error {
	
	return a.client.Call("UnregisterAgent", 0, agent).Store()
	
}

