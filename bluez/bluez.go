package bluez

import "github.com/godbus/dbus/introspect"

const (
	OrgBluezPath      = "/org/bluez"
	OrgBluezInterface = "org.bluez"

	//ObjectManagerInterface the dbus object manager interface
	ObjectManagerInterface = "org.freedesktop.DBus.ObjectManager"
	//InterfacesRemoved the DBus signal member for InterfacesRemoved
	InterfacesRemoved = "org.freedesktop.DBus.ObjectManager.InterfacesRemoved"
	//InterfacesAdded the DBus signal member for InterfacesAdded
	InterfacesAdded = "org.freedesktop.DBus.ObjectManager.InterfacesAdded"

	//PropertiesInterface the DBus properties interface
	PropertiesInterface = "org.freedesktop.DBus.Properties"
	//PropertiesChanged the DBus properties interface and member
	PropertiesChanged = "org.freedesktop.DBus.Properties.PropertiesChanged"

	// Introspectable introspectable interface
	Introspectable = "org.freedesktop.DBus.Introspectable"

	// //Device1Interface the bluez interface for Device1
	// Device1Interface = "org.bluez.Device1"
	// //Adapter1Interface the bluez interface for Adapter1
	// Adapter1Interface = "org.bluez.Adapter1"
	// //GattService1Interface the bluez interface for GattService1
	// GattService1Interface = "org.bluez.GattService1"
	// //GattCharacteristic1Interface the bluez interface for GattCharacteristic1
	// GattCharacteristic1Interface = "org.bluez.GattCharacteristic1"
	// //GattDescriptor1Interface the bluez interface for GattDescriptor1
	// GattDescriptor1Interface = "org.bluez.GattDescriptor1"
	// //LEAdvertisement1Interface the bluez interface for LEAdvertisement1
	// LEAdvertisement1Interface = "org.bluez.LEAdvertisement1"
	// // Agent1Interface the bluez interface for Agent1
	// Agent1Interface = "org.bluez.Agent1"
	//
	// // Media API
	// // Media1Interface the bluez interface for Media1
	// Media1Interface = "org.bluez.Media1"
	// // MediaControl1Interface the bluez interface for MediaControl1
	// MediaControl1Interface = "org.bluez.MediaControl1"
	// // MediaPlayer1Interface the bluez interface for MediaPlayer1
	// MediaPlayer1Interface = "org.bluez.MediaPlayer1"
	// // MediaFolder1Interface the bluez interface for MediaFolder1
	// MediaFolder1Interface = "org.bluez.MediaFolder1"
	// // MediaItem1Interface the bluez interface for MediaItem1
	// MediaItem1Interface = "org.bluez.MediaItem1"
	// // MediaEndpoint1Interface the bluez interface for MediaEndpoint1
	// MediaEndpoint1Interface = "org.bluez.MediaEndpoint1"
	// // MediaTransport1Interface the bluez interface for MediaTransport1
	// MediaTransport1Interface = "org.bluez.MediaTransport1"

)

// ObjectManagerIntrospectDataString introspect ObjectManager description
const ObjectManagerIntrospectDataString = `
<interface name="org.freedesktop.DBus.ObjectManager">
	<method name="GetManagedObjects">
		<arg name="objects" type="a{oa{sa{sv}}}" direction="out" />
	</method>
	<signal name="InterfacesAdded">
		<arg name="object" type="o"/>
		<arg name="interfaces" type="a{sa{sv}}"/>
	</signal>
	<signal name="InterfacesRemoved">
		<arg name="object" type="o"/>
		<arg name="interfaces" type="as"/>
	</signal>
</interface>`

// ObjectManagerIntrospectData introspect ObjectManager description
var ObjectManagerIntrospectData = introspect.Interface{
	Name: "org.freedesktop.DBus.ObjectManager",
	Methods: []introspect.Method{
		{
			Name: "GetManagedObjects",
			Args: []introspect.Arg{
				{
					Name:      "objects",
					Type:      "a{oa{sa{sv}}}",
					Direction: "out",
				},
			},
		},
	},
	Signals: []introspect.Signal{
		{
			Name: "InterfacesAdded",
			Args: []introspect.Arg{
				{
					Name: "object",
					Type: "o",
				},
				{
					Name: "interfaces",
					Type: "a{sa{sv}}",
				},
			},
		},
		{
			Name: "InterfacesRemoved",
			Args: []introspect.Arg{
				{
					Name: "object",
					Type: "o",
				},
				{
					Name: "interfaces",
					Type: "as",
				},
			},
		},
	},
}
