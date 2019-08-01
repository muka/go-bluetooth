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
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
)

var MediaFolder1Interface = "org.bluez.MediaFolder1"


// NewMediaFolder1 create a new instance of MediaFolder1
//
// Args:
// 	servicePath: unique name
// 	objectPath: freely definable
func NewMediaFolder1(servicePath string, objectPath string) (*MediaFolder1, error) {
	a := new(MediaFolder1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: MediaFolder1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaFolder1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}

// NewMediaFolder1Controller create a new instance of MediaFolder1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/playerX
func NewMediaFolder1Controller(objectPath string) (*MediaFolder1, error) {
	a := new(MediaFolder1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: MediaFolder1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaFolder1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// MediaFolder1 MediaFolder1 hierarchy

type MediaFolder1 struct {
	client     *bluez.Client
	Properties *MediaFolder1Properties
}

// MediaFolder1Properties contains the exposed properties of an interface
type MediaFolder1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// NumberOfItems Number of items in the folder
	NumberOfItems uint32

	// Name Folder name:
  // Possible values:
  // "/Filesystem/...": Filesystem scope
  // "/NowPlaying/...": NowPlaying scope
  // Note: /NowPlaying folder might not be listed if player
  // is stopped, folders created by Search are virtual so
  // once another Search is perform or the folder is
  // changed using ChangeFolder it will no longer be listed.
  // Filters
	Name string

	// Start Offset of the first item.
  // Default value: 0
	Start uint32

	// End Offset of the last item.
  // Default value: NumbeOfItems
	End uint32

	// Attributes Item properties that should be included in the list.
  // Possible Values:
  // "title", "artist", "album", "genre",
  // "number-of-tracks", "number", "duration"
  // Default Value: All
	Attributes []string

}

func (p *MediaFolder1Properties) Lock() {
	p.lock.Lock()
}

func (p *MediaFolder1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *MediaFolder1) Close() {
	a.client.Disconnect()
}


// ToMap convert a MediaFolder1Properties to map
func (a *MediaFolder1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MediaFolder1Properties
func (a *MediaFolder1Properties) FromMap(props map[string]interface{}) (*MediaFolder1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MediaFolder1Properties
func (a *MediaFolder1Properties) FromDBusMap(props map[string]dbus.Variant) (*MediaFolder1Properties, error) {
	s := new(MediaFolder1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MediaFolder1) GetProperties() (*MediaFolder1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MediaFolder1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MediaFolder1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *MediaFolder1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *MediaFolder1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Search Return a folder object containing the search result.
// To list the items found use the folder object returned
// and pass to ChangeFolder.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaFolder1) Search(value string, filter map[string]dbus.Variant) (dbus.ObjectPath, error) {
	
	var val0 dbus.ObjectPath
	err := a.client.Call("Search", 0, value, filter).Store(&val0)
	return val0, err	
}

//ChangeFolder Change current folder.
// Note: By changing folder the items of previous folder
// might be destroyed and have to be listed again, the
// exception is NowPlaying folder which should be always
// present while the player is active.
// Possible Errors: org.bluez.Error.InvalidArguments
// org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaFolder1) ChangeFolder(folder dbus.ObjectPath) error {
	
	return a.client.Call("ChangeFolder", 0, folder).Store()
	
}

