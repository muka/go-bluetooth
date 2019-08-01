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

package network



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Network1Interface = "org.bluez.Network1"


// NewNetwork1 create a new instance of Network1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewNetwork1(objectPath string) (*Network1, error) {
	a := new(Network1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Network1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Network1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Network1 Network hierarchy

type Network1 struct {
	client     *bluez.Client
	Properties *Network1Properties
}

// Network1Properties contains the exposed properties of an interface
type Network1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Interface Indicates the network interface name when available.
	Interface string

	// UUID Indicates the connection role when available.
	UUID string

	// Connected Indicates if the device is connected.
	Connected bool

}

func (p *Network1Properties) Lock() {
	p.lock.Lock()
}

func (p *Network1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Network1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Network1Properties to map
func (a *Network1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Network1Properties
func (a *Network1Properties) FromMap(props map[string]interface{}) (*Network1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Network1Properties
func (a *Network1Properties) FromDBusMap(props map[string]dbus.Variant) (*Network1Properties, error) {
	s := new(Network1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Network1) GetProperties() (*Network1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Network1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Network1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Network1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Network1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Connect Connect to the network device and return the network
// interface name. Examples of the interface name are
// bnep0, bnep1 etc.
// uuid can be either one of "gn", "panu" or "nap" (case
// insensitive) or a traditional string representation of
// UUID or a hexadecimal number.
// The connection will be closed and network device
func (a *Network1) Connect(uuid string) (string, error) {
	
	var val0 string
	err := a.client.Call("Connect", 0, uuid).Store(&val0)
	return val0, err	
}

//Disconnect Disconnect from the network device.
// To abort a connection attempt in case of errors or
// timeouts in the client it is fine to call this method.
// Possible errors: org.bluez.Error.Failed
func (a *Network1) Disconnect() error {
	
	return a.client.Call("Disconnect", 0, ).Store()
	
}

