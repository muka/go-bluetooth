// Code generated DO NOT EDIT

package mesh



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/muka/go-bluetooth/util"
   "github.com/muka/go-bluetooth/props"
   "github.com/godbus/dbus"
)

var Application1Interface = "org.bluez.mesh.Application1"


// NewApplication1 create a new instance of Application1
//
// Args:
// - servicePath: unique name
// - objectPath: <app_root>
func NewApplication1(servicePath string, objectPath dbus.ObjectPath) (*Application1, error) {
	a := new(Application1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  servicePath,
			Iface: Application1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Application1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
Application1 Mesh Application Hierarchy

*/
type Application1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Application1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// Application1Properties contains the exposed properties of an interface
type Application1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

	/*
	CompanyID A 16-bit Bluetooth-assigned Company Identifier of the vendor as
		defined by Bluetooth SIG
	*/
	CompanyID uint16

	/*
	ProductID A 16-bit vendor-assigned product identifier
	*/
	ProductID uint16

	/*
	VersionID A 16-bit vendor-assigned product version identifier
	*/
	VersionID uint16

	/*
	CRPL A 16-bit minimum number of replay protection list entries
	*/
	CRPL uint16

}

//Lock access to properties
func (p *Application1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *Application1Properties) Unlock() {
	p.lock.Unlock()
}




// SetCompanyID set CompanyID value
func (a *Application1) SetCompanyID(v uint16) error {
	return a.SetProperty("CompanyID", v)
}



// GetCompanyID get CompanyID value
func (a *Application1) GetCompanyID() (uint16, error) {
	v, err := a.GetProperty("CompanyID")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}




// SetProductID set ProductID value
func (a *Application1) SetProductID(v uint16) error {
	return a.SetProperty("ProductID", v)
}



// GetProductID get ProductID value
func (a *Application1) GetProductID() (uint16, error) {
	v, err := a.GetProperty("ProductID")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}




// SetVersionID set VersionID value
func (a *Application1) SetVersionID(v uint16) error {
	return a.SetProperty("VersionID", v)
}



// GetVersionID get VersionID value
func (a *Application1) GetVersionID() (uint16, error) {
	v, err := a.GetProperty("VersionID")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}




// SetCRPL set CRPL value
func (a *Application1) SetCRPL(v uint16) error {
	return a.SetProperty("CRPL", v)
}



// GetCRPL get CRPL value
func (a *Application1) GetCRPL() (uint16, error) {
	v, err := a.GetProperty("CRPL")
	if err != nil {
		return uint16(0), err
	}
	return v.Value().(uint16), nil
}



// Close the connection
func (a *Application1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return Application1 object path
func (a *Application1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return Application1 dbus client
func (a *Application1) Client() *bluez.Client {
	return a.client
}

// Interface return Application1 interface
func (a *Application1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Application1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a Application1Properties to map
func (a *Application1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an Application1Properties
func (a *Application1Properties) FromMap(props map[string]interface{}) (*Application1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Application1Properties
func (a *Application1Properties) FromDBusMap(props map[string]dbus.Variant) (*Application1Properties, error) {
	s := new(Application1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *Application1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *Application1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *Application1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *Application1) GetProperties() (*Application1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Application1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Application1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Application1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Application1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *Application1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *Application1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
JoinComplete 
		This method is called when the node provisioning initiated

*/
func (a *Application1) JoinComplete(token uint64) error {
	
	return a.client.Call("JoinComplete", 0, token).Store()
	
}

/*
Join 
		The token parameter serves as a unique identifier of the
		particular node. The token must be preserved by the application
		in order to authenticate itself to the mesh daemon and attach to

*/
func (a *Application1) Join() (by a, error) {
	
	var val0 by a
	err := a.client.Call("Join", 0, ).Store(&val0)
	return val0, err	
}

/*
Leave 

*/
func (a *Application1) Leave() error {
	
	return a.client.Call("Leave", 0, ).Store()
	
}

/*
JoinFailed 

*/
func (a *Application1) JoinFailed(reason string) error {
	
	return a.client.Call("JoinFailed", 0, reason).Store()
	
}

/*
Join 
		The reason parameter identifies the reason for provisioning
		failure. The defined values are: "timeout", "bad-pdu",
		"confirmation-failed", "out-of-resources", "decryption-error",
		"unexpected-error", "cannot-assign-addresses".

Properties:
	uint16 CompanyID [read-only]

		A 16-bit Bluetooth-assigned Company Identifier of the vendor as
		defined by Bluetooth SIG

	uint16 ProductID [read-only]

		A 16-bit vendor-assigned product identifier

	uint16 VersionID [read-only]

		A 16-bit vendor-assigned product version identifier

	uint16 CRPL [read-only, optional]

		A 16-bit minimum number of replay protection list entries



*/
func (a *Application1) Join() error {
	
	return a.client.Call("Join", 0, ).Store()
	
}

