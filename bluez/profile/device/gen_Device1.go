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

package device



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var Device1Interface = "org.bluez.Device1"


// NewDevice1 create a new instance of Device1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewDevice1(objectPath dbus.ObjectPath) (*Device1, error) {
	a := new(Device1)
	a.propertiesSignal = make(chan *dbus.Signal)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: Device1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Device1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// Device1 Device hierarchy

type Device1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	Properties 				*Device1Properties
}

// Device1Properties contains the exposed properties of an interface
type Device1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Address The Bluetooth device address of the remote device.
	Address string

	// Icon Proposed icon name according to the freedesktop.org
  // icon naming specification.
	Icon string

	// Alias The name alias for the remote device. The alias can
  // be used to have a different friendly name for the
  // remote device.
  // In case no alias is set, it will return the remote
  // device name. Setting an empty string as alias will
  // convert it back to the remote device name.
  // When resetting the alias with an empty string, the
  // property will default back to the remote name.
	Alias string

	// LegacyPairing Set to true if the device only supports the pre-2.1
  // pairing mechanism. This property is useful during
  // device discovery to anticipate whether legacy or
  // simple pairing will occur if pairing is initiated.
  // Note that this property can exhibit false-positives
  // in the case of Bluetooth 2.1 (or newer) devices that
  // have disabled Extended Inquiry Response support.
	LegacyPairing bool

	// RSSI Received Signal Strength Indicator of the remote
  // device (inquiry or advertising).
	RSSI int16

	// AdvertisingData The Advertising Data of the remote device. Keys are
  // are 8 bits AD Type followed by data as byte array.
  // Note: Only types considered safe to be handled by
  // application are exposed.
  // Possible values:
  // <type> <byte array>
  // ...
  // Example:
  // <Transport Discovery> <Organization Flags...>
  // 0x26                   0x01         0x01...
	AdvertisingData map[string]interface{}

	// AddressType The Bluetooth device Address Type. For dual-mode and
  // BR/EDR only devices this defaults to "public". Single
  // mode LE devices may have either value. If remote device
  // uses privacy than before pairing this represents address
  // type used for connection and Identity Address after
  // pairing.
  // Possible values:
  // "public" - Public address
  // "random" - Random address
	AddressType string

	// Name The Bluetooth remote name. This value can not be
  // changed. Use the Alias property instead.
  // This value is only present for completeness. It is
  // better to always use the Alias property when
  // displaying the devices name.
  // If the Alias property is unset, it will reflect
  // this value which makes it more convenient.
	Name string

	// Connected Indicates if the remote device is currently connected.
  // A PropertiesChanged signal indicate changes to this
  // status.
	Connected bool

	// Trusted Indicates if the remote is seen as trusted. This
  // setting can be changed by the application.
	Trusted bool

	// Blocked If set to true any incoming connections from the
  // device will be immediately rejected. Any device
  // drivers will also be removed and no new ones will
  // be probed as long as the device is blocked.
	Blocked bool

	// TxPower Advertised transmitted power level (inquiry or
  // advertising).
	TxPower int16

	// Paired Indicates if the remote device is paired.
	Paired bool

	// Adapter The object path of the adapter the device belongs to.
	Adapter dbus.ObjectPath

	// Modalias Remote Device ID information in modalias format
  // used by the kernel and udev.
	Modalias string

	// ManufacturerData Manufacturer specific advertisement data. Keys are
  // 16 bits Manufacturer ID followed by its byte array
  // value.
	ManufacturerData map[uint16]interface{}

	// ServiceData Service advertisement data. Keys are the UUIDs in
  // string format followed by its byte array value.
	ServiceData map[string]interface{}

	// AdvertisingFlags The Advertising Data Flags of the remote device.
	AdvertisingFlags []byte

	// Class The Bluetooth class of device of the remote device.
	Class uint32

	// Appearance External appearance of device, as found on GAP service.
	Appearance uint16

	// UUIDs List of 128-bit UUIDs that represents the available
  // remote services.
	UUIDs []string

	// ServicesResolved Indicate whether or not service discovery has been
  // resolved.
	ServicesResolved bool

}

func (p *Device1Properties) Lock() {
	p.lock.Lock()
}

func (p *Device1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *Device1) Close() {
	
	a.unregisterSignal()
	
	a.client.Disconnect()
}

// Path return Device1 object path
func (a *Device1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return Device1 interface
func (a *Device1) Interface() string {
	return a.client.Config.Iface
}


// ToMap convert a Device1Properties to map
func (a *Device1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an Device1Properties
func (a *Device1Properties) FromMap(props map[string]interface{}) (*Device1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Device1Properties
func (a *Device1Properties) FromDBusMap(props map[string]dbus.Variant) (*Device1Properties, error) {
	s := new(Device1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *Device1) GetProperties() (*Device1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Device1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Device1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Device1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Device1) unregisterSignal() {
	if a.propertiesSignal == nil {
		a.propertiesSignal <- nil
	}
}

// WatchProperties updates on property changes
func (a *Device1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *Device1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}





//Connect This is a generic method to connect any profiles
// the remote device supports that can be connected
// to and have been flagged as auto-connectable on
// our side. If only subset of profiles is already
// connected it will try to connect currently disconnected
// ones.
// If at least one profile was connected successfully this
// method will indicate success.
// For dual-mode devices only one bearer is connected at
// time, the conditions are in the following order:
// 1. Connect the disconnected bearer if already
// connected.
// 2. Connect first the bonded bearer. If no
// bearers are bonded or both are skip and check
// latest seen bearer.
// 3. Connect last seen bearer, in case the
// timestamps are the same BR/EDR takes
// precedence.
// Possible errors: org.bluez.Error.NotReady
// org.bluez.Error.Failed
// org.bluez.Error.InProgress
// org.bluez.Error.AlreadyConnected
func (a *Device1) Connect() error {
	
	return a.client.Call("Connect", 0, ).Store()
	
}

//Disconnect This method gracefully disconnects all connected
// profiles and then terminates low-level ACL connection.
// ACL connection will be terminated even if some profiles
// were not disconnected properly e.g. due to misbehaving
// device.
// This method can be also used to cancel a preceding
// Connect call before a reply to it has been received.
// For non-trusted devices connected over LE bearer calling
// this method will disable incoming connections until
// Connect method is called again.
// Possible errors: org.bluez.Error.NotConnected
func (a *Device1) Disconnect() error {
	
	return a.client.Call("Disconnect", 0, ).Store()
	
}

//ConnectProfile This method connects a specific profile of this
// device. The UUID provided is the remote service
// UUID for the profile.
// Possible errors: org.bluez.Error.Failed
// org.bluez.Error.InProgress
// org.bluez.Error.InvalidArguments
// org.bluez.Error.NotAvailable
// org.bluez.Error.NotReady
func (a *Device1) ConnectProfile(uuid string) error {
	
	return a.client.Call("ConnectProfile", 0, uuid).Store()
	
}

//DisconnectProfile This method disconnects a specific profile of
// this device. The profile needs to be registered
// client profile.
// There is no connection tracking for a profile, so
// as long as the profile is registered this will always
// succeed.
// Possible errors: org.bluez.Error.Failed
// org.bluez.Error.InProgress
// org.bluez.Error.InvalidArguments
// org.bluez.Error.NotSupported
func (a *Device1) DisconnectProfile(uuid string) error {
	
	return a.client.Call("DisconnectProfile", 0, uuid).Store()
	
}

//Pair This method will connect to the remote device,
func (a *Device1) Pair() error {
	
	return a.client.Call("Pair", 0, ).Store()
	
}

//CancelPairing This method can be used to cancel a pairing
// operation initiated by the Pair method.
// Possible errors: org.bluez.Error.DoesNotExist
// org.bluez.Error.Failed
func (a *Device1) CancelPairing() error {
	
	return a.client.Call("CancelPairing", 0, ).Store()
	
}

