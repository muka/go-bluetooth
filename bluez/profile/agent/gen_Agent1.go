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

package agent



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Agent1Interface = "org.bluez.Agent1"


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

//RequestPinCode This method gets called when the service daemon
// needs to get the passkey for an authentication.
// The return value should be a string of 1-16 characters
// length. The string can be alphanumeric.
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Agent1) RequestPinCode(device dbus.ObjectPath) (string, error) {
	
	var val0 string
	err := a.client.Call("RequestPinCode", 0, device).Store(&val0)
	return val0, err	
}

//DisplayPinCode This method gets called when the service daemon
// needs to display a pincode for an authentication.
// An empty reply should be returned. When the pincode
// needs no longer to be displayed, the Cancel method
// of the agent will be called.
// This is used during the pairing process of keyboards
// that don't support Bluetooth 2.1 Secure Simple Pairing,
// in contrast to DisplayPasskey which is used for those
// that do.
// This method will only ever be called once since
// older keyboards do not support typing notification.
// Note that the PIN will always be a 6-digit number,
// zero-padded to 6 digits. This is for harmony with
// the later specification.
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Agent1) DisplayPinCode(device dbus.ObjectPath, pincode string) error {
	
	return a.client.Call("DisplayPinCode", 0, device, pincode).Store()
	
}

//RequestPasskey This method gets called when the service daemon
// needs to get the passkey for an authentication.
// The return value should be a numeric value
// between 0-999999.
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Agent1) RequestPasskey(device dbus.ObjectPath) (uint32, error) {
	
	var val0 uint32
	err := a.client.Call("RequestPasskey", 0, device).Store(&val0)
	return val0, err	
}

//DisplayPasskey This method gets called when the service daemon
// needs to display a passkey for an authentication.
// The entered parameter indicates the number of already
// typed keys on the remote side.
// An empty reply should be returned. When the passkey
// needs no longer to be displayed, the Cancel method
// of the agent will be called.
// During the pairing process this method might be
// called multiple times to update the entered value.
// Note that the passkey will always be a 6-digit number,
// so the display should be zero-padded at the start if
// the value contains less than 6 digits.
func (a *Agent1) DisplayPasskey(device dbus.ObjectPath, passkey uint32, entered uint16) error {
	
	return a.client.Call("DisplayPasskey", 0, device, passkey, entered).Store()
	
}

//RequestConfirmation This method gets called when the service daemon
// needs to confirm a passkey for an authentication.
// To confirm the value it should return an empty reply
// or an error in case the passkey is invalid.
// Note that the passkey will always be a 6-digit number,
// so the display should be zero-padded at the start if
// the value contains less than 6 digits.
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Agent1) RequestConfirmation(device dbus.ObjectPath, passkey uint32) error {
	
	return a.client.Call("RequestConfirmation", 0, device, passkey).Store()
	
}

//RequestAuthorization This method gets called to request the user to
// authorize an incoming pairing attempt which
// would in other circumstances trigger the just-works
// model, or when the user plugged in a device that
// implements cable pairing. In the latter case, the
// device would not be connected to the adapter via
// Bluetooth yet.
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Agent1) RequestAuthorization(device dbus.ObjectPath) error {
	
	return a.client.Call("RequestAuthorization", 0, device).Store()
	
}

//AuthorizeService This method gets called when the service daemon
// needs to authorize a connection/service request.
// Possible errors: org.bluez.Error.Rejected
// org.bluez.Error.Canceled
func (a *Agent1) AuthorizeService(device dbus.ObjectPath, uuid string) error {
	
	return a.client.Call("AuthorizeService", 0, device, uuid).Store()
	
}

//Cancel This method gets called to indicate that the agent
// request failed before a reply was returned.
func (a *Agent1) Cancel() error {
	
	return a.client.Call("Cancel", 0, ).Store()
	
}

