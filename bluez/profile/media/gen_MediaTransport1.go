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
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var MediaTransport1Interface = "org.bluez.MediaTransport1"


// NewMediaTransport1 create a new instance of MediaTransport1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/fdX
func NewMediaTransport1(objectPath dbus.ObjectPath) (*MediaTransport1, error) {
	a := new(MediaTransport1)
	a.propertiesSignal = make(chan *dbus.Signal)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: MediaTransport1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaTransport1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// MediaTransport1 MediaTransport1 hierarchy

type MediaTransport1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	Properties 				*MediaTransport1Properties
}

// MediaTransport1Properties contains the exposed properties of an interface
type MediaTransport1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Device Device object which the transport is connected to.
	Device dbus.ObjectPath

	// UUID UUID of the profile which the transport is for.
	UUID string

	// Codec Assigned number of codec that the transport support.
  // The values should match the profile specification which
  // is indicated by the UUID.
	Codec byte

	// Configuration Configuration blob, it is used as it is so the size and
  // byte order must match.
	Configuration []byte

	// State Indicates the state of the transport. Possible
  // values are:
  // "idle": not streaming
  // "pending": streaming but not acquired
  // "active": streaming and acquired
	State string

	// Delay Optional. Transport delay in 1/10 of millisecond, this
  // property is only writeable when the transport was
  // acquired by the sender.
	Delay uint16

	// Volume Optional. Indicates volume level of the transport,
  // this property is only writeable when the transport was
  // acquired by the sender.
  // Possible Values: 0-127
	Volume uint16

}

func (p *MediaTransport1Properties) Lock() {
	p.lock.Lock()
}

func (p *MediaTransport1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *MediaTransport1) Close() {
	
	a.unregisterSignal()
	
	a.client.Disconnect()
}

// Path return MediaTransport1 object path
func (a *MediaTransport1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return MediaTransport1 interface
func (a *MediaTransport1) Interface() string {
	return a.client.Config.Iface
}


// ToMap convert a MediaTransport1Properties to map
func (a *MediaTransport1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MediaTransport1Properties
func (a *MediaTransport1Properties) FromMap(props map[string]interface{}) (*MediaTransport1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MediaTransport1Properties
func (a *MediaTransport1Properties) FromDBusMap(props map[string]dbus.Variant) (*MediaTransport1Properties, error) {
	s := new(MediaTransport1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MediaTransport1) GetProperties() (*MediaTransport1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MediaTransport1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MediaTransport1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *MediaTransport1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *MediaTransport1) unregisterSignal() {
	if a.propertiesSignal == nil {
		a.propertiesSignal <- nil
	}
}

// WatchProperties updates on property changes
func (a *MediaTransport1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *MediaTransport1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}





//Acquire Acquire transport file descriptor and the MTU for read
// and write respectively.
// Possible Errors: org.bluez.Error.NotAuthorized
// org.bluez.Error.Failed
func (a *MediaTransport1) Acquire() (dbus.UnixFD, uint16, uint16, error) {
	
	var val0 dbus.UnixFD
  var val1 uint16
  var val2 uint16
	err := a.client.Call("Acquire", 0, ).Store(&val0, &val1, &val2)
	return val0, val1, val2, err	
}

//TryAcquire Acquire transport file descriptor only if the transport
// is in "pending" state at the time the message is
// received by BlueZ. Otherwise no request will be sent
// to the remote device and the function will just fail
// with org.bluez.Error.NotAvailable.
// Possible Errors: org.bluez.Error.NotAuthorized
// org.bluez.Error.Failed
// org.bluez.Error.NotAvailable
func (a *MediaTransport1) TryAcquire() (dbus.UnixFD, uint16, uint16, error) {
	
	var val0 dbus.UnixFD
  var val1 uint16
  var val2 uint16
	err := a.client.Call("TryAcquire", 0, ).Store(&val0, &val1, &val2)
	return val0, val1, val2, err	
}

//Release Releases file descriptor.
func (a *MediaTransport1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

