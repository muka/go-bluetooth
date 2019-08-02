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

var Message1Interface = "org.bluez.obex.Message1"


// NewMessage1 create a new instance of Message1
//
// Args:
// 	objectPath: [Session object path]/{message0,...}
func NewMessage1(objectPath string) (*Message1, error) {
	a := new(Message1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: Message1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Message1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Message1 Message hierarchy

type Message1 struct {
	client     *bluez.Client
	Properties *Message1Properties
}

// Message1Properties contains the exposed properties of an interface
type Message1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// RecipientAddress Message recipient address
	RecipientAddress string

	// Status Message reception status
  // Possible values: "complete",
  // "fractioned" and "notification"
	Status string

	// Deleted Message deleted flag
	Deleted bool

	// Sender Message sender name
	Sender string

	// ReplyTo Message Reply-To address
	ReplyTo string

	// Folder Folder which the message belongs to
	Folder string

	// Subject Message subject
	Subject string

	// Timestamp Message timestamp
	Timestamp string

	// SenderAddress Message sender address
	SenderAddress string

	// Read Message read flag
	Read bool

	// Recipient Message recipient name
	Recipient string

	// Type Message type
  // Possible values: "email", "sms-gsm",
  // "sms-cdma" and "mms"
  // uint64 Size [readonly]
  // Message size in bytes
	Type string

	// Priority Message priority flag
	Priority bool

	// Sent Message sent flag
	Sent bool

	// Protected Message protected flag
	Protected bool

}

func (p *Message1Properties) Lock() {
	p.lock.Lock()
}

func (p *Message1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Message1) Close() {
	a.client.Disconnect()
}


// ToMap convert a Message1Properties to map
func (a *Message1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Message1Properties
func (a *Message1Properties) FromMap(props map[string]interface{}) (*Message1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Message1Properties
func (a *Message1Properties) FromDBusMap(props map[string]dbus.Variant) (*Message1Properties, error) {
	s := new(Message1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Message1) GetProperties() (*Message1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Message1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Message1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *Message1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *Message1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Get Download message and store it in the target file.
// If an empty target file is given, a temporary file
// will be automatically generated.
// The returned path represents the newly created transfer,
// which should be used to find out if the content has been
// successfully transferred or if the operation fails.
// The properties of this transfer are also returned along
// with the object path, to avoid a call to GetProperties.
// Possible errors: org.bluez.obex.Error.InvalidArguments
// org.bluez.obex.Error.Failed
func (a *Message1) Get(targetfile string, attachment bool) (dbus.ObjectPath, map[string]dbus.Variant, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]dbus.Variant
	err := a.client.Call("Get", 0, targetfile, attachment).Store(&val0, &val1)
	return val0, val1, err	
}

