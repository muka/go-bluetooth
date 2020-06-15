package bluez

import "github.com/godbus/dbus/v5/introspect"

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
