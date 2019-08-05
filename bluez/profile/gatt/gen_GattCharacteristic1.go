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

package gatt



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var GattCharacteristic1Interface = "org.bluez.GattCharacteristic1"


// NewGattCharacteristic1 create a new instance of GattCharacteristic1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/serviceXX/charYYYY
func NewGattCharacteristic1(objectPath dbus.ObjectPath) (*GattCharacteristic1, error) {
	a := new(GattCharacteristic1)
	a.propertiesSignal = make(chan *dbus.Signal)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: GattCharacteristic1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(GattCharacteristic1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// GattCharacteristic1 Characteristic hierarchy
// For local GATT defined services, the object paths need to follow the service
// path hierarchy and are freely definable.
type GattCharacteristic1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	Properties 				*GattCharacteristic1Properties
}

// GattCharacteristic1Properties contains the exposed properties of an interface
type GattCharacteristic1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Descriptors 
	Descriptors []dbus.ObjectPath

	// UUID 128-bit characteristic UUID.
	UUID string

	// Service Object path of the GATT service the characteristic
  // belongs to.
	Service dbus.ObjectPath

	// Value The cached value of the characteristic. This property
  // gets updated only after a successful read request and
  // when a notification or indication is received, upon
  // which a PropertiesChanged signal will be emitted.
	Value []byte `dbus:"emit"`

	// WriteAcquired True, if this characteristic has been acquired by any
  // client using AcquireWrite.
  // For client properties is ommited in case
  // 'write-without-response' flag is not set.
  // For server the presence of this property indicates
  // that AcquireWrite is supported.
	WriteAcquired bool

	// NotifyAcquired True, if this characteristic has been acquired by any
  // client using AcquireNotify.
  // For client this properties is ommited in case 'notify'
  // flag is not set.
  // For server the presence of this property indicates
  // that AcquireNotify is supported.
	NotifyAcquired bool

	// Notifying True, if notifications or indications on this
  // characteristic are currently enabled.
	Notifying bool

	// Flags Defines how the characteristic value can be used. See
  // Core spec "Table 3.5: Characteristic Properties bit
  // field", and "Table 3.8: Characteristic Extended
  // Properties bit field". Allowed values:
  // "broadcast"
  // "read"
  // "write-without-response"
  // "write"
  // "notify"
  // "indicate"
  // "authenticated-signed-writes"
  // "reliable-write"
  // "writable-auxiliaries"
  // "encrypt-read"
  // "encrypt-write"
  // "encrypt-authenticated-read"
  // "encrypt-authenticated-write"
  // "secure-read" (Server only)
  // "secure-write" (Server only)
  // "authorize"
	Flags []string

}

func (p *GattCharacteristic1Properties) Lock() {
	p.lock.Lock()
}

func (p *GattCharacteristic1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *GattCharacteristic1) Close() {
	
	a.unregisterSignal()
	
	a.client.Disconnect()
}

// Path return GattCharacteristic1 object path
func (a *GattCharacteristic1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return GattCharacteristic1 interface
func (a *GattCharacteristic1) Interface() string {
	return a.client.Config.Iface
}


// ToMap convert a GattCharacteristic1Properties to map
func (a *GattCharacteristic1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an GattCharacteristic1Properties
func (a *GattCharacteristic1Properties) FromMap(props map[string]interface{}) (*GattCharacteristic1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an GattCharacteristic1Properties
func (a *GattCharacteristic1Properties) FromDBusMap(props map[string]dbus.Variant) (*GattCharacteristic1Properties, error) {
	s := new(GattCharacteristic1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *GattCharacteristic1) GetProperties() (*GattCharacteristic1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *GattCharacteristic1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *GattCharacteristic1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *GattCharacteristic1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *GattCharacteristic1) unregisterSignal() {
	if a.propertiesSignal == nil {
		a.propertiesSignal <- nil
	}
}

// WatchProperties updates on property changes
func (a *GattCharacteristic1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *GattCharacteristic1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}





//ReadValue Issues a request to read the value of the
// characteristic and returns the value if the
// operation was successful.
// Possible options: "offset": uint16 offset
// "device": Object Device (Server only)
// Possible Errors: org.bluez.Error.Failed
// org.bluez.Error.InProgress
// org.bluez.Error.NotPermitted
// org.bluez.Error.NotAuthorized
// org.bluez.Error.InvalidOffset
// org.bluez.Error.NotSupported
func (a *GattCharacteristic1) ReadValue(options map[string]interface{}) ([]byte, error) {
	
	var val0 []byte
	err := a.client.Call("ReadValue", 0, options).Store(&val0)
	return val0, err	
}

//WriteValue Issues a request to write the value of the
// characteristic.
// Possible options: "offset": Start offset
// "device": Device path (Server only)
// "link": Link type (Server only)
// "prepare-authorize": boolean Is prepare
// authorization
// request
// Possible Errors: org.bluez.Error.Failed
// org.bluez.Error.InProgress
// org.bluez.Error.NotPermitted
// org.bluez.Error.InvalidValueLength
// org.bluez.Error.NotAuthorized
// org.bluez.Error.NotSupported
func (a *GattCharacteristic1) WriteValue(value []byte, options map[string]interface{}) error {
	
	return a.client.Call("WriteValue", 0, value, options).Store()
	
}

//AcquireWrite Acquire file descriptor and MTU for writing. Usage of
// WriteValue will be locked causing it to return
// NotPermitted error.
// For server the MTU returned shall be equal or smaller
// than the negotiated MTU.
// For client it only works with characteristic that has
// WriteAcquired property which relies on
// write-without-response Flag.
// To release the lock the client shall close the file
// descriptor, a HUP is generated in case the device
// is disconnected.
// Note: the MTU can only be negotiated once and is
// symmetric therefore this method may be delayed in
// order to have the exchange MTU completed, because of
// that the file descriptor is closed during
// reconnections as the MTU has to be renegotiated.
// Possible options: "device": Object Device (Server only)
// "MTU": Exchanged MTU (Server only)
// "link": Link type (Server only)
// Possible Errors: org.bluez.Error.Failed
// org.bluez.Error.NotSupported
func (a *GattCharacteristic1) AcquireWrite(options map[string]interface{}) (dbus.UnixFD, uint16, error) {
	
	var val0 dbus.UnixFD
  var val1 uint16
	err := a.client.Call("AcquireWrite", 0, options).Store(&val0, &val1)
	return val0, val1, err	
}

//AcquireNotify Acquire file descriptor and MTU for notify. Usage of
// StartNotify will be locked causing it to return
// NotPermitted error.
// For server the MTU returned shall be equal or smaller
// than the negotiated MTU.
// Only works with characteristic that has NotifyAcquired
// which relies on notify Flag and no other client have
// called StartNotify.
// Notification are enabled during this procedure so
// StartNotify shall not be called, any notification
// will be dispatched via file descriptor therefore the
// Value property is not affected during the time where
// notify has been acquired.
// To release the lock the client shall close the file
// descriptor, a HUP is generated in case the device
// is disconnected.
// Note: the MTU can only be negotiated once and is
// symmetric therefore this method may be delayed in
// order to have the exchange MTU completed, because of
// that the file descriptor is closed during
// reconnections as the MTU has to be renegotiated.
// Possible options: "device": Object Device (Server only)
// "MTU": Exchanged MTU (Server only)
// "link": Link type (Server only)
// Possible Errors: org.bluez.Error.Failed
// org.bluez.Error.NotSupported
func (a *GattCharacteristic1) AcquireNotify(options map[string]interface{}) (dbus.UnixFD, uint16, error) {
	
	var val0 dbus.UnixFD
  var val1 uint16
	err := a.client.Call("AcquireNotify", 0, options).Store(&val0, &val1)
	return val0, val1, err	
}

//StartNotify Starts a notification session from this characteristic
// if it supports value notifications or indications.
// Possible Errors: org.bluez.Error.Failed
// org.bluez.Error.NotPermitted
// org.bluez.Error.InProgress
// org.bluez.Error.NotSupported
func (a *GattCharacteristic1) StartNotify() error {
	
	return a.client.Call("StartNotify", 0, ).Store()
	
}

//StopNotify This method will cancel any previous StartNotify
// transaction. Note that notifications from a
// characteristic are shared between sessions thus
// calling StopNotify will release a single session.
// Possible Errors: org.bluez.Error.Failed
func (a *GattCharacteristic1) StopNotify() error {
	
	return a.client.Call("StopNotify", 0, ).Store()
	
}

//Confirm This method doesn't expect a reply so it is just a
// confirmation that value was received.
// Possible Errors: org.bluez.Error.Failed
func (a *GattCharacteristic1) Confirm() error {
	
	return a.client.Call("Confirm", 0, ).Store()
	
}

