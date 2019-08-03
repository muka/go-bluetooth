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

var MessageAccess1Interface = "org.bluez.obex.MessageAccess1"


// NewMessageAccess1 create a new instance of MessageAccess1
//
// Args:
// 	objectPath: [Session object path]
func NewMessageAccess1(objectPath string) (*MessageAccess1, error) {
	a := new(MessageAccess1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: MessageAccess1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MessageAccess1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// MessageAccess1 Message Access hierarchy

type MessageAccess1 struct {
	client     *bluez.Client
	Properties *MessageAccess1Properties
}

// MessageAccess1Properties contains the exposed properties of an interface
type MessageAccess1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

func (p *MessageAccess1Properties) Lock() {
	p.lock.Lock()
}

func (p *MessageAccess1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *MessageAccess1) Close() {
	a.client.Disconnect()
}


// ToMap convert a MessageAccess1Properties to map
func (a *MessageAccess1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MessageAccess1Properties
func (a *MessageAccess1Properties) FromMap(props map[string]interface{}) (*MessageAccess1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MessageAccess1Properties
func (a *MessageAccess1Properties) FromDBusMap(props map[string]dbus.Variant) (*MessageAccess1Properties, error) {
	s := new(MessageAccess1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MessageAccess1) GetProperties() (*MessageAccess1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MessageAccess1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MessageAccess1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *MessageAccess1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *MessageAccess1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//SetFolder Set working directory for current session, *name* may
// be the directory name or '..[/dir]'.
// Possible errors: org.bluez.obex.Error.InvalidArguments
// org.bluez.obex.Error.Failed
func (a *MessageAccess1) SetFolder(name string) error {
	
	return a.client.Call("SetFolder", 0, name).Store()
	
}

//ListFolders Returns a dictionary containing information about
// the current folder content.
// The following keys are defined:
// string Name : Folder name
// Possible filters: Offset and MaxCount
// Possible errors: org.bluez.obex.Error.InvalidArguments
// org.bluez.obex.Error.Failed
func (a *MessageAccess1) ListFolders(filter map[string]interface{}) ([]map[string]interface{}, error) {
	
	var val0 []map[string]interface{}
	err := a.client.Call("ListFolders", 0, filter).Store(&val0)
	return val0, err	
}

//ListFilterFields Return all available fields that can be used in Fields
// filter.
// Possible errors: None
func (a *MessageAccess1) ListFilterFields() ([]string, error) {
	
	var val0 []string
	err := a.client.Call("ListFilterFields", 0, ).Store(&val0)
	return val0, err	
}

//ListMessages Returns an array containing the messages found in the
// given subfolder of the current folder, or in the
// current folder if folder is empty.
// Possible Filters: Offset, MaxCount, SubjectLength, Fields,
// Type, PeriodStart, PeriodEnd, Status, Recipient, Sender,
// Priority
// Each message is represented by an object path followed
// by a dictionary of the properties.
// Properties:
// string Subject:
// Message subject
// string Timestamp:
// Message timestamp
// string Sender:
// Message sender name
// string SenderAddress:
// Message sender address
// string ReplyTo:
// Message Reply-To address
// string Recipient:
// Message recipient name
// string RecipientAddress:
// Message recipient address
// string Type:
// Message type
// Possible values: "email", "sms-gsm",
// "sms-cdma" and "mms"
// uint64 Size:
// Message size in bytes
// boolean Text:
// Message text flag
// Specifies whether message has textual
// content or is binary only
// string Status:
// Message status
// Possible values for received messages:
// "complete", "fractioned", "notification"
// Possible values for sent messages:
// "delivery-success", "sending-success",
// "delivery-failure", "sending-failure"
// uint64 AttachmentSize:
// Message overall attachment size in bytes
// boolean Priority:
// Message priority flag
// boolean Read:
// Message read flag
// boolean Sent:
// Message sent flag
// boolean Protected:
// Message protected flag
// Possible errors: org.bluez.obex.Error.InvalidArguments
// org.bluez.obex.Error.Failed
func (a *MessageAccess1) ListMessages(folder string, filter map[string]interface{}) ([]dbus.ObjectPath, map[string]interface{}, error) {
	
	var val0 []dbus.ObjectPath
  var val1 map[string]interface{}
	err := a.client.Call("ListMessages", 0, folder, filter).Store(&val0, &val1)
	return val0, val1, err	
}

//UpdateInbox 
func (a *MessageAccess1) UpdateInbox() error {
	
	return a.client.Call("UpdateInbox", 0, ).Store()
	
}

//PushMessage Transfer a message (in bMessage format) to the
// remote device.
// The message is transferred either to the given
// subfolder of the current folder, or to the current
// folder if folder is empty.
// Possible args: Transparent, Retry, Charset
// The returned path represents the newly created transfer,
// which should be used to find out if the content has been
// successfully transferred or if the operation fails.
// The properties of this transfer are also returned along
// with the object path, to avoid a call to GetAll.
// Possible errors: org.bluez.obex.Error.InvalidArguments
// org.bluez.obex.Error.Failed
// Filter:		uint16 Offset:
// Offset of the first item, default is 0
// uint16 MaxCount:
// Maximum number of items, default is 1024
// byte SubjectLength:
// Maximum length of the Subject property in the
// message, default is 256
// array{string} Fields:
// Message fields, default is all values.
// Possible values can be query with ListFilterFields.
// array{string} Types:
// Filter messages by type.
// Possible values: "sms", "email", "mms".
// string PeriodBegin:
// Filter messages by starting period.
// Possible values: Date in "YYYYMMDDTHHMMSS" format.
// string PeriodEnd:
// Filter messages by ending period.
// Possible values: Date in "YYYYMMDDTHHMMSS" format.
// boolean Read:
// Filter messages by read flag.
// Possible values: True for read or False for unread
// string Recipient:
// Filter messages by recipient address.
// string Sender:
// Filter messages by sender address.
// boolean Priority:
// Filter messages by priority flag.
// Possible values: True for high priority or False for
// non-high priority
func (a *MessageAccess1) PushMessage(sourcefile string, folder string, args map[string]interface{}) error {
	
	return a.client.Call("PushMessage", 0, sourcefile, folder, args).Store()
	
}

