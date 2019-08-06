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
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Message1Interface = "org.bluez.obex.Message1"


// NewMessage1 create a new instance of Message1
//
// Args:
// 	objectPath: [Session object path]/{message0,...}
func NewMessage1(objectPath dbus.ObjectPath) (*Message1, error) {
	a := new(Message1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.obex",
			Iface: Message1Interface,
			Path:  dbus.ObjectPath(objectPath),
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
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Message1Properties
}

// Message1Properties contains the exposed properties of an interface
type Message1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Read Message read flag
	Read bool

	// Protected Message protected flag
	Protected bool

	// Folder Folder which the message belongs to
	Folder string

	// Subject Message subject
	Subject string

	// RecipientAddress Message recipient address
	RecipientAddress string

	// Priority Message priority flag
	Priority bool

	// Deleted Message deleted flag
	Deleted bool

	// Sent Message sent flag
	Sent bool

	// ReplyTo Message Reply-To address
	ReplyTo string

	// Recipient Message recipient name
	Recipient string

	// Status Message reception status
  // Possible values: "complete",
  // "fractioned" and "notification"
	Status string

	// Type Message type
  // Possible values: "email", "sms-gsm",
  // "sms-cdma" and "mms"
  // uint64 Size [readonly]
  // Message size in bytes
	Type string

	// Timestamp Message timestamp
	Timestamp string

	// Sender Message sender name
	Sender string

	// SenderAddress Message sender address
	SenderAddress string

}

func (p *Message1Properties) Lock() {
	p.lock.Lock()
}

func (p *Message1Properties) Unlock() {
	p.lock.Unlock()
}


// SetRead set Read value
func (a *Message1) SetRead(v bool) error {
	return a.SetProperty("Read", v)
}

// GetRead get Read value
func (a *Message1) GetRead() (bool, error) {
	v, err := a.GetProperty("Read")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetProtected set Protected value
func (a *Message1) SetProtected(v bool) error {
	return a.SetProperty("Protected", v)
}

// GetProtected get Protected value
func (a *Message1) GetProtected() (bool, error) {
	v, err := a.GetProperty("Protected")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetFolder set Folder value
func (a *Message1) SetFolder(v string) error {
	return a.SetProperty("Folder", v)
}

// GetFolder get Folder value
func (a *Message1) GetFolder() (string, error) {
	v, err := a.GetProperty("Folder")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetSubject set Subject value
func (a *Message1) SetSubject(v string) error {
	return a.SetProperty("Subject", v)
}

// GetSubject get Subject value
func (a *Message1) GetSubject() (string, error) {
	v, err := a.GetProperty("Subject")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetRecipientAddress set RecipientAddress value
func (a *Message1) SetRecipientAddress(v string) error {
	return a.SetProperty("RecipientAddress", v)
}

// GetRecipientAddress get RecipientAddress value
func (a *Message1) GetRecipientAddress() (string, error) {
	v, err := a.GetProperty("RecipientAddress")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetPriority set Priority value
func (a *Message1) SetPriority(v bool) error {
	return a.SetProperty("Priority", v)
}

// GetPriority get Priority value
func (a *Message1) GetPriority() (bool, error) {
	v, err := a.GetProperty("Priority")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetDeleted set Deleted value
func (a *Message1) SetDeleted(v bool) error {
	return a.SetProperty("Deleted", v)
}

// GetDeleted get Deleted value
func (a *Message1) GetDeleted() (bool, error) {
	v, err := a.GetProperty("Deleted")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetSent set Sent value
func (a *Message1) SetSent(v bool) error {
	return a.SetProperty("Sent", v)
}

// GetSent get Sent value
func (a *Message1) GetSent() (bool, error) {
	v, err := a.GetProperty("Sent")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}

// SetReplyTo set ReplyTo value
func (a *Message1) SetReplyTo(v string) error {
	return a.SetProperty("ReplyTo", v)
}

// GetReplyTo get ReplyTo value
func (a *Message1) GetReplyTo() (string, error) {
	v, err := a.GetProperty("ReplyTo")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetRecipient set Recipient value
func (a *Message1) SetRecipient(v string) error {
	return a.SetProperty("Recipient", v)
}

// GetRecipient get Recipient value
func (a *Message1) GetRecipient() (string, error) {
	v, err := a.GetProperty("Recipient")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetStatus set Status value
func (a *Message1) SetStatus(v string) error {
	return a.SetProperty("Status", v)
}

// GetStatus get Status value
func (a *Message1) GetStatus() (string, error) {
	v, err := a.GetProperty("Status")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetType set Type value
func (a *Message1) SetType(v string) error {
	return a.SetProperty("Type", v)
}

// GetType get Type value
func (a *Message1) GetType() (string, error) {
	v, err := a.GetProperty("Type")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetTimestamp set Timestamp value
func (a *Message1) SetTimestamp(v string) error {
	return a.SetProperty("Timestamp", v)
}

// GetTimestamp get Timestamp value
func (a *Message1) GetTimestamp() (string, error) {
	v, err := a.GetProperty("Timestamp")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetSender set Sender value
func (a *Message1) SetSender(v string) error {
	return a.SetProperty("Sender", v)
}

// GetSender get Sender value
func (a *Message1) GetSender() (string, error) {
	v, err := a.GetProperty("Sender")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}

// SetSenderAddress set SenderAddress value
func (a *Message1) SetSenderAddress(v string) error {
	return a.SetProperty("SenderAddress", v)
}

// GetSenderAddress get SenderAddress value
func (a *Message1) GetSenderAddress() (string, error) {
	v, err := a.GetProperty("SenderAddress")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}


// Close the connection
func (a *Message1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return Message1 object path
func (a *Message1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return Message1 interface
func (a *Message1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Message1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

	if a.objectManagerSignal == nil {
		if a.objectManager == nil {
			om, err := bluez.GetObjectManager()
			if err != nil {
				return nil, nil, err
			}
			a.objectManager = om
		}

		s, err := a.objectManager.Register()
		if err != nil {
			return nil, nil, err
		}
		a.objectManagerSignal = s
	}

	cancel := func() {
		if a.objectManagerSignal == nil {
			return
		}
		a.objectManagerSignal <- nil
		a.objectManager.Unregister(a.objectManagerSignal)
		a.objectManagerSignal = nil
	}

	return a.objectManagerSignal, cancel, nil
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

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Message1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Message1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *Message1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

				// updates [*]Properties struct when a property change
				s := reflect.ValueOf(a.Properties).Elem()
				// exported field
				f := s.FieldByName(field)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						x := reflect.ValueOf(val.Value())
						a.Properties.Lock()
						f.Set(x)
						a.Properties.Unlock()
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

func (a *Message1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
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
func (a *Message1) Get(targetfile string, attachment bool) (dbus.ObjectPath, map[string]interface{}, error) {
	
	var val0 dbus.ObjectPath
  var val1 map[string]interface{}
	err := a.client.Call("Get", 0, targetfile, attachment).Store(&val0, &val1)
	return val0, val1, err	
}

