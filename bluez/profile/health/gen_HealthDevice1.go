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
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var HealthDevice1Interface = "org.bluez.HealthDevice1"


// NewHealthDevice1 create a new instance of HealthDevice1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewHealthDevice1(objectPath string) (*HealthDevice1, error) {
	a := new(HealthDevice1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: HealthDevice1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(HealthDevice1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// HealthDevice1 HealthDevice hierarchy

type HealthDevice1 struct {
	client     *bluez.Client
	Properties *HealthDevice1Properties
}

// HealthDevice1Properties contains the exposed properties of an interface
type HealthDevice1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// MainChannel The first reliable channel opened. It is needed by
  // upper applications in order to send specific protocol
  // data units. The first reliable can change after a
  // reconnection.
	MainChannel dbus.ObjectPath

}

func (p *HealthDevice1Properties) Lock() {
	p.lock.Lock()
}

func (p *HealthDevice1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *HealthDevice1) Close() {
	a.client.Disconnect()
}


// ToMap convert a HealthDevice1Properties to map
func (a *HealthDevice1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an HealthDevice1Properties
func (a *HealthDevice1Properties) FromMap(props map[string]interface{}) (*HealthDevice1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an HealthDevice1Properties
func (a *HealthDevice1Properties) FromDBusMap(props map[string]dbus.Variant) (*HealthDevice1Properties, error) {
	s := new(HealthDevice1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *HealthDevice1) GetProperties() (*HealthDevice1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *HealthDevice1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *HealthDevice1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *HealthDevice1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *HealthDevice1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Echo Sends an echo petition to the remote service. Returns
// True if response matches with the buffer sent. If some
// error is detected False value is returned.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.OutOfRange
func (a *HealthDevice1) Echo() (bool, error) {
	
	var val0 bool
	err := a.client.Call("Echo", 0, ).Store(&val0)
	return val0, err	
}

//CreateChannel Creates a new data channel.  The configuration should
// indicate the channel quality of service using one of
// this values "reliable", "streaming", "any".
// Returns the object path that identifies the data
// channel that is already connected.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.HealthError
func (a *HealthDevice1) CreateChannel(application dbus.ObjectPath, configuration string) (dbus.ObjectPath, error) {
	
	var val0 dbus.ObjectPath
	err := a.client.Call("CreateChannel", 0, application, configuration).Store(&val0)
	return val0, err	
}

//DestroyChannel Destroys the data channel object. Only the creator of
// the channel or the creator of the HealthApplication
// that received the data channel will be able to destroy
// it.
// Possible errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.NotFound
// org.bluez.Error.NotAllowed
func (a *HealthDevice1) DestroyChannel(channel dbus.ObjectPath) error {
	
	return a.client.Call("DestroyChannel", 0, channel).Store()
	
}

