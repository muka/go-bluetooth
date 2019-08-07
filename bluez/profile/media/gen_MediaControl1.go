package media



import (
  "sync"
  "github.com/muka/go-bluetooth/bluez"
  "reflect"
  "github.com/fatih/structs"
  "github.com/muka/go-bluetooth/util"
  "github.com/godbus/dbus"
  "fmt"
)

var MediaControl1Interface = "org.bluez.MediaControl1"


// NewMediaControl1 create a new instance of MediaControl1
//
// Args:
// - objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX
func NewMediaControl1(objectPath dbus.ObjectPath) (*MediaControl1, error) {
	a := new(MediaControl1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: MediaControl1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaControl1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}

// NewMediaControl1FromAdapterID create a new instance of MediaControl1
// adapterID: ID of an adapter eg. hci0
func NewMediaControl1FromAdapterID(adapterID string) (*MediaControl1, error) {
	a := new(MediaControl1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez",
			Iface: MediaControl1Interface,
			Path:  dbus.ObjectPath(fmt.Sprintf("/org/bluez/%s", adapterID)),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(MediaControl1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
MediaControl1 Media Control hierarchy

*/
type MediaControl1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*MediaControl1Properties
}

// MediaControl1Properties contains the exposed properties of an interface
type MediaControl1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	/*
	Player Addressed Player object path.
	*/
	Player dbus.ObjectPath

	/*
	Connected 
	*/
	Connected bool

}

//Lock access to properties
func (p *MediaControl1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *MediaControl1Properties) Unlock() {
	p.lock.Unlock()
}


// SetPlayer set Player value
func (a *MediaControl1) SetPlayer(v dbus.ObjectPath) error {
	return a.SetProperty("Player", v)
}

// GetPlayer get Player value
func (a *MediaControl1) GetPlayer() (dbus.ObjectPath, error) {
	v, err := a.GetProperty("Player")
	if err != nil {
		return dbus.ObjectPath(""), err
	}
	return v.Value().(dbus.ObjectPath), nil
}

// SetConnected set Connected value
func (a *MediaControl1) SetConnected(v bool) error {
	return a.SetProperty("Connected", v)
}

// GetConnected get Connected value
func (a *MediaControl1) GetConnected() (bool, error) {
	v, err := a.GetProperty("Connected")
	if err != nil {
		return false, err
	}
	return v.Value().(bool), nil
}


// Close the connection
func (a *MediaControl1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return MediaControl1 object path
func (a *MediaControl1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Interface return MediaControl1 interface
func (a *MediaControl1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *MediaControl1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a MediaControl1Properties to map
func (a *MediaControl1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(a), nil
}

// FromMap convert a map to an MediaControl1Properties
func (a *MediaControl1Properties) FromMap(props map[string]interface{}) (*MediaControl1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an MediaControl1Properties
func (a *MediaControl1Properties) FromDBusMap(props map[string]dbus.Variant) (*MediaControl1Properties, error) {
	s := new(MediaControl1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// GetProperties load all available properties
func (a *MediaControl1) GetProperties() (*MediaControl1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *MediaControl1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *MediaControl1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *MediaControl1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *MediaControl1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *MediaControl1) WatchProperties() (chan *bluez.PropertyChanged, error) {

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

func (a *MediaControl1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




/*
Play 
			Resume playback.


*/
func (a *MediaControl1) Play() error {
	
	return a.client.Call("Play", 0, ).Store()
	
}

/*
Pause 
			Pause playback.


*/
func (a *MediaControl1) Pause() error {
	
	return a.client.Call("Pause", 0, ).Store()
	
}

/*
Stop 
			Stop playback.


*/
func (a *MediaControl1) Stop() error {
	
	return a.client.Call("Stop", 0, ).Store()
	
}

/*
Next 
			Next item.


*/
func (a *MediaControl1) Next() error {
	
	return a.client.Call("Next", 0, ).Store()
	
}

/*
Previous 
			Previous item.


*/
func (a *MediaControl1) Previous() error {
	
	return a.client.Call("Previous", 0, ).Store()
	
}

/*
VolumeUp 
			Adjust remote volume one step up


*/
func (a *MediaControl1) VolumeUp() error {
	
	return a.client.Call("VolumeUp", 0, ).Store()
	
}

/*
VolumeDown 
			Adjust remote volume one step down


*/
func (a *MediaControl1) VolumeDown() error {
	
	return a.client.Call("VolumeDown", 0, ).Store()
	
}

/*
FastForward 
			Fast forward playback, this action is only stopped
			when another method in this interface is called.


*/
func (a *MediaControl1) FastForward() error {
	
	return a.client.Call("FastForward", 0, ).Store()
	
}

/*
Rewind 
			Rewind playback, this action is only stopped
			when another method in this interface is called.

Properties

		boolean Connected [readonly]

		object Player [readonly, optional]

			Addressed Player object path.



*/
func (a *MediaControl1) Rewind() error {
	
	return a.client.Call("Rewind", 0, ).Store()
	
}

