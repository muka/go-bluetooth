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

package advertising



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var LEAdvertisement1Interface = "org.bluez.LEAdvertisement1"


// NewLEAdvertisement1 create a new instance of LEAdvertisement1
//
// Args:
// 	objectPath: freely definable
func NewLEAdvertisement1(objectPath string) (*LEAdvertisement1, error) {
	a := new(LEAdvertisement1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: LEAdvertisement1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(LEAdvertisement1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// LEAdvertisement1 LE Advertisement Data hierarchy
// Specifies the Advertisement Data to be broadcast and some advertising
// parameters.  Properties which are not present will not be included in the
// data.  Required advertisement data types will always be included.
// All UUIDs are 128-bit versions in the API, and 16 or 32-bit
// versions of the same UUID will be used in the advertising data as appropriate.
type LEAdvertisement1 struct {
	client     *bluez.Client
	Properties *LEAdvertisement1Properties
}

// LEAdvertisement1Properties contains the exposed properties of an interface
type LEAdvertisement1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Timeout Timeout of the advertisement in seconds. This defines
  // the lifetime of the advertisement.
	Timeout uint16

	// Type Determines the type of advertising packet requested.
  // Possible values: "broadcast" or "peripheral"
	Type string

	// ManufacturerData Manufactuer Data fields to include in
  // the Advertising Data.  Keys are the Manufacturer ID
  // to associate with the data.
	ManufacturerData map[uint16]interface{}

	// DiscoverableTimeout The discoverable timeout in seconds. A value of zero
  // means that the timeout is disabled and it will stay in
  // discoverable/limited mode forever.
  // Note: This property shall not be set when Type is set
  // to broadcast.
	DiscoverableTimeout uint16

	// Includes List of features to be included in the advertising
  // packet.
  // Possible values: as found on
  // LEAdvertisingManager.SupportedIncludes
	Includes []string

	// Appearance Appearance to be used in the advertising report.
  // Possible values: as found on GAP Service.
	Appearance uint16

	// Duration Duration of the advertisement in seconds. If there are
  // other applications advertising no duration is set the
  // default is 2 seconds.
	Duration uint16

	// ServiceUUIDs List of UUIDs to include in the "Service UUID" field of
  // the Advertising Data.
	ServiceUUIDs []string

	// SolicitUUIDs Array of UUIDs to include in "Service Solicitation"
  // Advertisement Data.
	SolicitUUIDs []string

	// ServiceData Service Data elements to include. The keys are the
  // UUID to associate with the data.
	ServiceData map[string]interface{}

	// Data Advertising Type to include in the Advertising
  // Data. Key is the advertising type and value is the
  // data as byte array.
  // Note: Types already handled by other properties shall
  // not be used.
  // Possible values:
  // <type> <byte array>
  // ...
  // Example:
  // <Transport Discovery> <Organization Flags...>
  // 0x26                   0x01         0x01...
	Data map[byte]interface{}

	// Discoverable Advertise as general discoverable. When present this
  // will override adapter Discoverable property.
  // Note: This property shall not be set when Type is set
  // to broadcast.
	Discoverable bool

	// LocalName Local name to be used in the advertising report. If the
  // string is too big to fit into the packet it will be
  // truncated.
  // If this property is available 'local-name' cannot be
  // present in the Includes.
	LocalName string

}

func (p *LEAdvertisement1Properties) Lock() {
	p.lock.Lock()
}

func (p *LEAdvertisement1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *LEAdvertisement1) Close() {
	a.client.Disconnect()
}


// ToMap convert a LEAdvertisement1Properties to map
func (a *LEAdvertisement1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an LEAdvertisement1Properties
func (a *LEAdvertisement1Properties) FromMap(props map[string]interface{}) (*LEAdvertisement1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an LEAdvertisement1Properties
func (a *LEAdvertisement1Properties) FromDBusMap(props map[string]dbus.Variant) (*LEAdvertisement1Properties, error) {
	s := new(LEAdvertisement1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *LEAdvertisement1) GetProperties() (*LEAdvertisement1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *LEAdvertisement1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *LEAdvertisement1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *LEAdvertisement1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *LEAdvertisement1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Release This method gets called when the service daemon
// removes the Advertisement. A client can use it to do
// cleanup tasks. There is no need to call
// UnregisterAdvertisement because when this method gets
// called it has already been unregistered.
func (a *LEAdvertisement1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

