package media



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
  log "github.com/sirupsen/logrus"
   "reflect"
   "github.com/fatih/structs"
   "github.com/muka/go-bluetooth/util"
   "github.com/godbus/dbus"
)

var MediaTransport1Interface = "org.bluez.MediaTransport1"


// NewMediaTransport1 create a new instance of MediaTransport1
//
// Args:
// - objectPath: [variable prefix]/{hci0,hci1,...}/dev_XX_XX_XX_XX_XX_XX/fdX
func NewMediaTransport1(objectPath dbus.ObjectPath) (*MediaTransport1, error) {
	a := new(MediaTransport1)
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


/*
MediaTransport1 MediaTransport1 hierarchy

*/
type MediaTransport1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*MediaTransport1Properties
}

// MediaTransport1Properties contains the exposed properties of an interface
type MediaTransport1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	/*
	Volume Optional. Indicates volume level of the transport,
			this property is only writeable when the transport was
			acquired by the sender.

			Possible Values: 0-127
	*/
	Volume uint16

	/*
	Device Device object which the transport is connected to.
	*/
	Device dbus.ObjectPath

	/*
	UUID UUID of the profile which the transport is for.
	*/
	UUID string

	/*
	Codec Assigned number of codec that the transport support.
			The values should match the profile specification which
			is indicated by the UUID.
	*/
	Codec byte

	/*
	Configuration Configuration blob, it is used as it is so the size and
			byte order must match.
	*/
	Configuration []byte

	/*
	State Indicates the state of the transport. Possible
			values are:
				"idle": not streaming
				"pending": streaming but not acquired
				"active": streaming and acquired
	*/
	State string

	/*
	Delay Optional. Transport delay in 1/10 of millisecond, this
			property is only writeable when the transport was
			acquired by the sender.
	*/
	Delay uint16

}

//Lock access to properties
func (p *MediaTransport1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *MediaTransport1Properties) Unlock() {
	p.lock.Unlock()
}




// SetVolume set Volume value
func (a *MediaTransport1) SetVolume(v uint16) error {
	return a.SetProperty("Volume", v)
}



// GetVolume get Volume value
func (a *MediaTransport1) GetVolume() (uint16, error) {
	v, err := a.GetProperty("Volume")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}






// GetDevice get Device value
func (a *MediaTransport1) GetDevice() (dbus.ObjectPath, error) {
	v, err := a.GetProperty("Device")
	if err != nil {
		return dbus.ObjectPath(""), err
	}
	return v.Value().(dbus.ObjectPath), nil
}






// GetUUID get UUID value
func (a *MediaTransport1) GetUUID() (string, error) {
	v, err := a.GetProperty("UUID")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}






// GetCodec get Codec value
func (a *MediaTransport1) GetCodec() (byte, error) {
	v, err := a.GetProperty("Codec")
	if err != nil {
		return byte(0), err
	}
	return v.Value().(byte), nil
}






// GetConfiguration get Configuration value
func (a *MediaTransport1) GetConfiguration() ([]byte, error) {
	v, err := a.GetProperty("Configuration")
	if err != nil {
		return []byte{}, err
	}
	return v.Value().([]byte), nil
}






// GetState get State value
func (a *MediaTransport1) GetState() (string, error) {
	v, err := a.GetProperty("State")
	if err != nil {
		return "", err
	}
	return v.Value().(string), nil
}




// SetDelay set Delay value
func (a *MediaTransport1) SetDelay(v uint16) error {
	return a.SetProperty("Delay", v)
}



// GetDelay get Delay value
func (a *MediaTransport1) GetDelay() (uint16, error) {
	v, err := a.GetProperty("Delay")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}



// Close the connection
func (a *MediaTransport1) Close() {
	
	a.unregisterPropertiesSignal()
	
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

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *MediaTransport1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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
func (a *MediaTransport1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *MediaTransport1) WatchProperties() (chan *bluez.PropertyChanged, error) {

	// channel, err := a.client.Register(a.Path(), a.Interface())
	channel, err := a.client.Register(a.Path(), bluez.PropertiesInterface)
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
						// map[*]variant -> map[*]interface{}
						ok, err := util.AssignMapVariantToInterface(f, x)
						if err != nil {
							log.Errorf("Failed to set %s: %s", f.String(), err)
							continue
						}
						// direct assignment
						if !ok {
							f.Set(x)
						}
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

func (a *MediaTransport1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	ch <- nil
	close(ch)
	return nil
}




/*
Acquire 
			Acquire transport file descriptor and the MTU for read
			and write respectively.

			Possible Errors: org.bluez.Error.NotAuthorized
					 org.bluez.Error.Failed


*/
func (a *MediaTransport1) Acquire() (dbus.UnixFD, uint16, uint16, error) {
	
	var val0 dbus.UnixFD
  var val1 uint16
  var val2 uint16
	err := a.client.Call("Acquire", 0, ).Store(&val0, &val1, &val2)
	return val0, val1, val2, err	
}

/*
TryAcquire 
			Acquire transport file descriptor only if the transport
			is in "pending" state at the time the message is
			received by BlueZ. Otherwise no request will be sent
			to the remote device and the function will just fail
			with org.bluez.Error.NotAvailable.

			Possible Errors: org.bluez.Error.NotAuthorized
					 org.bluez.Error.Failed
					 org.bluez.Error.NotAvailable


*/
func (a *MediaTransport1) TryAcquire() (dbus.UnixFD, uint16, uint16, error) {
	
	var val0 dbus.UnixFD
  var val1 uint16
  var val2 uint16
	err := a.client.Call("TryAcquire", 0, ).Store(&val0, &val1, &val2)
	return val0, val1, val2, err	
}

/*
Release 
			Releases file descriptor.


*/
func (a *MediaTransport1) Release() error {
	
	return a.client.Call("Release", 0, ).Store()
	
}

