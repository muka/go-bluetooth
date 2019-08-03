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

var MediaPlayer1Interface = "org.bluez.MediaPlayer1"


// NewMediaPlayer1 create a new instance of MediaPlayer1
//
// Args:
// 	objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/playerX
func NewMediaPlayer1(objectPath string) (*MediaPlayer1, error) {
	a := new(MediaPlayer1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: MediaPlayer1Interface,
			Path:  objectPath,
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaPlayer1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


// MediaPlayer1 MediaPlayer1 hierarchy

type MediaPlayer1 struct {
	client     *bluez.Client
	Properties *MediaPlayer1Properties
}

// MediaPlayer1Properties contains the exposed properties of an interface
type MediaPlayer1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	// Position Playback position in milliseconds. Changing the
  // position may generate additional events that will be
  // sent to the remote device. When position is 0 it means
  // the track is starting and when it's greater than or
  // equal to track's duration the track has ended. Note
  // that even if duration is not available in metadata it's
  // possible to signal its end by setting position to the
  // maximum uint32 value.
	Position uint32

	// Title Track title name
	Title string

	// TrackNumber Track number
	TrackNumber uint32

	// Duration Track duration in milliseconds
	Duration uint32

	// Genre Track genre name
	Genre string

	// NumberOfTracks Number of tracks in total
	NumberOfTracks uint32

	// Name Player name
	Name string

	// Equalizer Possible values: "off" or "on"
	Equalizer string

	// Repeat Possible values: "off", "singletrack", "alltracks" or
  // "group"
	Repeat string

	// Scan Possible values: "off", "alltracks" or "group"
	Scan string

	// Track Track metadata.
  // Possible values:
	Track map[string]interface{}

	// Album Track album name
	Album string

	// Device Device object path.
	Device dbus.ObjectPath

	// Playlist Playlist object path.
	Playlist dbus.ObjectPath

	// Browsable If present indicates the player can be browsed using
  // MediaFolder interface.
  // Possible values:
  // True: Supported and active
  // False: Supported but inactive
  // Note: If supported but inactive clients can enable it
  // by using MediaFolder interface but it might interfere
  // in the playback of other players.
	Browsable bool

	// Searchable If present indicates the player can be searched using
  // MediaFolder interface.
  // Possible values:
  // True: Supported and active
  // False: Supported but inactive
  // Note: If supported but inactive clients can enable it
  // by using MediaFolder interface but it might interfere
  // in the playback of other players.
	Searchable bool

	// Shuffle Possible values: "off", "alltracks" or "group"
	Shuffle string

	// Status Possible status: "playing", "stopped", "paused",
  // "forward-seek", "reverse-seek"
  // or "error"
	Status string

	// Artist Track artist name
	Artist string

	// Type Player type
  // Possible values:
  // "Audio"
  // "Video"
  // "Audio Broadcasting"
  // "Video Broadcasting"
	Type string

	// Subtype Player subtype
  // Possible values:
  // "Audio Book"
  // "Podcast"
	Subtype string

}

func (p *MediaPlayer1Properties) Lock() {
	p.lock.Lock()
}

func (p *MediaPlayer1Properties) Unlock() {
	p.lock.Unlock()
}

// Close the connection
func (a *MediaPlayer1) Close() {
	a.client.Disconnect()
}


// ToMap convert a MediaPlayer1Properties to map
func (a *MediaPlayer1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MediaPlayer1Properties
func (a *MediaPlayer1Properties) FromMap(props map[string]interface{}) (*MediaPlayer1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MediaPlayer1Properties
func (a *MediaPlayer1Properties) FromDBusMap(props map[string]dbus.Variant) (*MediaPlayer1Properties, error) {
	s := new(MediaPlayer1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MediaPlayer1) GetProperties() (*MediaPlayer1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MediaPlayer1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MediaPlayer1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// Register for changes signalling
func (a *MediaPlayer1) Register() (chan *dbus.Signal, error) {
	return a.client.Register(a.client.Config.Path, bluez.PropertiesInterface)
}

// Unregister for changes signalling
func (a *MediaPlayer1) Unregister(signal chan *dbus.Signal) error {
	return a.client.Unregister(a.client.Config.Path, bluez.PropertiesInterface, signal)
}



//Play Resume playback.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaPlayer1) Play() error {
	
	return a.client.Call("Play", 0, ).Store()
	
}

//Pause Pause playback.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaPlayer1) Pause() error {
	
	return a.client.Call("Pause", 0, ).Store()
	
}

//Stop Stop playback.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaPlayer1) Stop() error {
	
	return a.client.Call("Stop", 0, ).Store()
	
}

//Next Next item.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaPlayer1) Next() error {
	
	return a.client.Call("Next", 0, ).Store()
	
}

//Previous Previous item.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaPlayer1) Previous() error {
	
	return a.client.Call("Previous", 0, ).Store()
	
}

//FastForward Fast forward playback, this action is only stopped
// when another method in this interface is called.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaPlayer1) FastForward() error {
	
	return a.client.Call("FastForward", 0, ).Store()
	
}

//Rewind Rewind playback, this action is only stopped
// when another method in this interface is called.
// Possible Errors: org.bluez.Error.NotSupported
// org.bluez.Error.Failed
func (a *MediaPlayer1) Rewind() error {
	
	return a.client.Call("Rewind", 0, ).Store()
	
}

