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
  "reflect"
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
	a.propertiesSignal = make(chan *dbus.Signal)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: AgentManager1Interface,
			Path:  dbus.ObjectPath("/org/bluez/obex"),
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
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	Properties 				*AgentManager1Properties
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
	
	a.unregisterSignal()
	
	a.client.Disconnect()
}

// Path return AgentManager1 object path
func (a *AgentManager1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return AgentManager1 interface
func (a *AgentManager1) Interface() string {
	return a.client.Config.Iface
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

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *AgentManager1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *AgentManager1) unregisterSignal() {
	if a.propertiesSignal == nil {
		a.propertiesSignal <- nil
	}
}

// WatchProperties updates on property changes
func (a *AgentManager1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *AgentManager1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
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

