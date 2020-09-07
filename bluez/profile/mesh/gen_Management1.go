// Code generated DO NOT EDIT

package mesh



import (
   "sync"
   "github.com/muka/go-bluetooth/bluez"
   "github.com/muka/go-bluetooth/util"
   "github.com/muka/go-bluetooth/props"
   "github.com/godbus/dbus/v5"
)

var Management1Interface = "org.bluez.mesh.Management1"


// NewManagement1 create a new instance of Management1
//
// Args:

func NewManagement1(objectPath dbus.ObjectPath) (*Management1, error) {
	a := new(Management1)
	a.client = bluez.NewClient(
		&bluez.Config{
			Name:  "org.bluez.mesh",
			Iface: Management1Interface,
			Path:  dbus.ObjectPath(objectPath),
			Bus:   bluez.SystemBus,
		},
	)
	
	a.Properties = new(Management1Properties)

	_, err := a.GetProperties()
	if err != nil {
		return nil, err
	}
	
	return a, nil
}


/*
Management1 Mesh Provisioning Hierarchy

*/
type Management1 struct {
	client     				*bluez.Client
	propertiesSignal 	chan *dbus.Signal
	objectManagerSignal chan *dbus.Signal
	objectManager       *bluez.ObjectManager
	Properties 				*Management1Properties
	watchPropertiesChannel chan *dbus.Signal
}

// Management1Properties contains the exposed properties of an interface
type Management1Properties struct {
	lock sync.RWMutex `dbus:"ignore"`

}

//Lock access to properties
func (p *Management1Properties) Lock() {
	p.lock.Lock()
}

//Unlock access to properties
func (p *Management1Properties) Unlock() {
	p.lock.Unlock()
}



// Close the connection
func (a *Management1) Close() {
	
	a.unregisterPropertiesSignal()
	
	a.client.Disconnect()
}

// Path return Management1 object path
func (a *Management1) Path() dbus.ObjectPath {
	return a.client.Config.Path
}

// Client return Management1 dbus client
func (a *Management1) Client() *bluez.Client {
	return a.client
}

// Interface return Management1 interface
func (a *Management1) Interface() string {
	return a.client.Config.Iface
}

// GetObjectManagerSignal return a channel for receiving updates from the ObjectManager
func (a *Management1) GetObjectManagerSignal() (chan *dbus.Signal, func(), error) {

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


// ToMap convert a Management1Properties to map
func (a *Management1Properties) ToMap() (map[string]interface{}, error) {
	return props.ToMap(a), nil
}

// FromMap convert a map to an Management1Properties
func (a *Management1Properties) FromMap(props map[string]interface{}) (*Management1Properties, error) {
	props1 := map[string]dbus.Variant{}
	for k, val := range props {
		props1[k] = dbus.MakeVariant(val)
	}
	return a.FromDBusMap(props1)
}

// FromDBusMap convert a map to an Management1Properties
func (a *Management1Properties) FromDBusMap(props map[string]dbus.Variant) (*Management1Properties, error) {
	s := new(Management1Properties)
	err := util.MapToStruct(s, props)
	return s, err
}

// ToProps return the properties interface
func (a *Management1) ToProps() bluez.Properties {
	return a.Properties
}

// GetWatchPropertiesChannel return the dbus channel to receive properties interface
func (a *Management1) GetWatchPropertiesChannel() chan *dbus.Signal {
	return a.watchPropertiesChannel
}

// SetWatchPropertiesChannel set the dbus channel to receive properties interface
func (a *Management1) SetWatchPropertiesChannel(c chan *dbus.Signal) {
	a.watchPropertiesChannel = c
}

// GetProperties load all available properties
func (a *Management1) GetProperties() (*Management1Properties, error) {
	a.Properties.Lock()
	err := a.client.GetProperties(a.Properties)
	a.Properties.Unlock()
	return a.Properties, err
}

// SetProperty set a property
func (a *Management1) SetProperty(name string, value interface{}) error {
	return a.client.SetProperty(name, value)
}

// GetProperty get a property
func (a *Management1) GetProperty(name string) (dbus.Variant, error) {
	return a.client.GetProperty(name)
}

// GetPropertiesSignal return a channel for receiving udpdates on property changes
func (a *Management1) GetPropertiesSignal() (chan *dbus.Signal, error) {

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
func (a *Management1) unregisterPropertiesSignal() {
	if a.propertiesSignal != nil {
		a.propertiesSignal <- nil
		a.propertiesSignal = nil
	}
}

// WatchProperties updates on property changes
func (a *Management1) WatchProperties() (chan *bluez.PropertyChanged, error) {
	return bluez.WatchProperties(a)
}

func (a *Management1) UnwatchProperties(ch chan *bluez.PropertyChanged) error {
	return bluez.UnwatchProperties(a, ch)
}




/*
UnprovisionedScan 		This method is used by the application that supports
		org.bluez.mesh.Provisioner1 interface to start listening
		(scanning) for unprovisioned devices in the area. Scanning
		will continue for the specified number of seconds, or, if 0 is
		specified, then continuously until UnprovisionedScanCancel() is
		called or if AddNode() method is called.
		Each time a unique unprovisioned beacon is heard, the
		ScanResult() method on the app will be called with the result.
		PossibleErrors:
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.NotAuthorized
			org.bluez.mesh.Error.Busy

*/
func (a *Management1) UnprovisionedScan(seconds uint16) error {
	
	return a.client.Call("UnprovisionedScan", 0, seconds).Store()
	
}

/*
UnprovisionedScanCancel 
*/
func (a *Management1) UnprovisionedScanCancel() error {
	
	return a.client.Call("UnprovisionedScanCancel", 0, ).Store()
	
}

/*
AddNode 		This method is used by the application that supports
		org.bluez.mesh.Provisioner1 interface to add the
		unprovisioned device specified by uuid, to the Network.
		The uuid parameter is a 16-byte array that contains Device UUID
		of the unprovisioned device to be added to the network.
		PossibleErrors:
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.NotAuthorized

*/
func (a *Management1) AddNode(uuid []byte) error {
	
	return a.client.Call("AddNode", 0, uuid).Store()
	
}

/*
CreateSubnet 		This method is used by the application to generate and add a new
		network subnet key.
		The net_index parameter is a 12-bit value (0x001-0xFFF)
		specifying which net key to add.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.AlreadyExists

*/
func (a *Management1) CreateSubnet(net_index uint16) error {
	
	return a.client.Call("CreateSubnet", 0, net_index).Store()
	
}

/*
ImportSubnet 		This method is used by the application to add a network subnet
		key, that was originally generated by a remote Config Client.
		The net_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which net key to add.
		The net_key parameter is the 16-byte value of the net key being
		imported.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.AlreadyExists

*/
func (a *Management1) ImportSubnet(net_index uint16, net_key []byte) error {
	
	return a.client.Call("ImportSubnet", 0, net_index, net_key).Store()
	
}

/*
UpdateSubnet 		This method is used by the application to generate a new network
		subnet key, and set it's key refresh state to Phase 1.
		The net_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which net key to update. Note that the subnet must
		exist prior to updating.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.DoesNotExist
			org.bluez.mesh.Error.Busy

*/
func (a *Management1) UpdateSubnet(net_index uint16) error {
	
	return a.client.Call("UpdateSubnet", 0, net_index).Store()
	
}

/*
DeleteSubnet 		This method is used by the application that to delete a subnet.
		The net_index parameter is a 12-bit value (0x001-0xFFF)
		specifying which net key to delete. The primary net key (0x000)
		may not be deleted.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.InvalidArguments

*/
func (a *Management1) DeleteSubnet(net_index uint16) error {
	
	return a.client.Call("DeleteSubnet", 0, net_index).Store()
	
}

/*
SetKeyPhase 		This method is used to set the master key update phase of the
		given subnet. When finalizing the procedure, it is important
		to CompleteAppKeyUpdate() on all app keys that have been
		updated during the procedure prior to setting phase 3.
		The net_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which subnet phase to set.
		The phase parameter is used to cycle the local key database
		through the phases as defined by the Mesh Profile Specification.
		Allowed values:
			0 - Cancel Key Refresh (May only be called from Phase 1,
				and should never be called once the new key has
				started propagating)
			1 - Invalid Argument (see NetKeyUpdate method)
			2 - Go to Phase 2 (May only be called from Phase 1)
			3 - Complete Key Refresh procedure (May only be called
				from Phase 2)
		This call affects the local bluetooth-meshd key database only.
		It is the responsibility of the application to maintain the key
		refresh phases per the Mesh Profile Specification.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.DoesNotExist

*/
func (a *Management1) SetKeyPhase(net_index uint16, phase uint8) error {
	
	return a.client.Call("SetKeyPhase", 0, net_index, phase).Store()
	
}

/*
CreateAppKey 		This method is used by the application to generate and add a new
		application key.
		The net_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which net key to bind the application key to.
		The app_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which app key to add.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.AlreadyExists
			org.bluez.mesh.Error.DoesNotExist

*/
func (a *Management1) CreateAppKey(net_index uint16, app_index uint16) error {
	
	return a.client.Call("CreateAppKey", 0, net_index, app_index).Store()
	
}

/*
ImportAppKey 		This method is used by the application to add an application
		key, that was originally generated by a remote Config Client.
		The net_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which net key to bind the application key to.
		The app_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which app key to import.
		The app_key parameter is the 16-byte value of the key being
		imported.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.AlreadyExists
			org.bluez.mesh.Error.DoesNotExist

*/
func (a *Management1) ImportAppKey(net_index uint16, app_index uint16, app_key []byte) error {
	
	return a.client.Call("ImportAppKey", 0, net_index, app_index, app_key).Store()
	
}

/*
UpdateAppKey 		This method is used by the application to generate a new
		application key.
		The app_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which app key to update. Note that the subnet that
		the key is bound to must exist and be in Phase 1.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments
			org.bluez.mesh.Error.DoesNotExist
			org.bluez.mesh.Error.Busy

*/
func (a *Management1) UpdateAppKey(app_index uint16) error {
	
	return a.client.Call("UpdateAppKey", 0, app_index).Store()
	
}

/*
DeleteAppKey 		This method is used by the application to delete an application
		key.
		The app_index parameter is a 12-bit value (0x000-0xFFF)
		specifying which app key to delete.
		This call affects the local bluetooth-meshd key database only.
		PossibleErrors:
			org.bluez.mesh.Error.InvalidArguments

*/
func (a *Management1) DeleteAppKey(app_index uint16) error {
	
	return a.client.Call("DeleteAppKey", 0, app_index).Store()
	
}

/*
ImportRemoteNode 		This method is used by the application to import a remote node
		that has been provisioned by an external process.
		The primary parameter specifies the unicast address of the
		the node being imported.
		The count parameter specifies the number of elements that are
		assigned to this remote node.
		The device_key parameter is the access layer key that will be
		will used to decrypt privledged messages from this remote node.
		This call affects the local bluetooth-meshd key database only.
		It is an error to call this with address range overlapping
		with local element addresses.
		PossibleErrors:
			org.bluez.mesh.Error.Failed
			org.bluez.mesh.Error.InvalidArguments

*/
func (a *Management1) ImportRemoteNode(primary uint16, count uint8, device_key []byte) error {
	
	return a.client.Call("ImportRemoteNode", 0, primary, count, device_key).Store()
	
}

/*
DeleteRemoteNode 		This method is used by the application to delete a remote node
		from the local device key database.
		The primary parameter specifies the unicast address of the
		the node being deleted.
		The count parameter specifies the number of elements that were
		assigned to the remote node.
		This call affects the local bluetooth-meshd key database only.
		It is an error to call this with address range overlapping
		with local element addresses.
		PossibleErrors:
			org.bluez.mesh.Error.InvalidArguments

*/
func (a *Management1) DeleteRemoteNode(primary uint16, count uint8) error {
	
	return a.client.Call("DeleteRemoteNode", 0, primary, count).Store()
	
}
