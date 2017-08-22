package bluez

const (

	//Device1Interface the bluez interface for Device1
	Device1Interface = "org.bluez.Device1"
	//Adapter1Interface the bluez interface for Adapter1
	Adapter1Interface = "org.bluez.Adapter1"
	//GattService1Interface the bluez interface for GattService1
	GattService1Interface = "org.bluez.GattService1"
	//GattCharacteristic1Interface the bluez interface for GattCharacteristic1
	GattCharacteristic1Interface = "org.bluez.GattCharacteristic1"
	//GattDescriptor1Interface the bluez interface for GattDescriptor1
	GattDescriptor1Interface = "org.bluez.GattDescriptor1"

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
)

// ObjectManagerIntrospectData introspect ObjectManager description
const ObjectManagerIntrospectData = `
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
