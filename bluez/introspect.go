package bluez

import "github.com/godbus/dbus/v5/introspect"

// GattService1IntrospectDataString interface definition
const GattService1IntrospectDataString = `
<interface name="org.bluez.GattService1">
  <property name="UUID" type="s" access="read"></property>
  <property name="Device" type="o" access="read"></property>
  <property name="Primary" type="b" access="read"></property>
  <property name="Characteristics" type="ao" access="read"></property>
</interface>
`

// GattDescriptor1IntrospectDataString interface definition
const GattDescriptor1IntrospectDataString = `
<interface name="org.bluez.GattDescriptor1">
  <method name="ReadValue">
    <arg name="value" type="ay" direction="out"/>
  </method>
  <method name="WriteValue">
    <arg name="value" type="ay" direction="in"/>
  </method>
  <property name="UUID" type="s" access="read"></property>
  <property name="Characteristic" type="o" access="read"></property>
  <property name="Value" type="ay" access="read"></property>
</interface>
`

//GattCharacteristic1IntrospectDataString interface definition
const GattCharacteristic1IntrospectDataString = `
<interface name="org.bluez.GattCharacteristic1">
  <method name="ReadValue">
    <arg name="value" type="ay" direction="out"/>
  </method>
  <method name="WriteValue">
    <arg name="value" type="ay" direction="in"/>
  </method>
  <method name="StartNotify"></method>
  <method name="StopNotify"></method>
  <property name="UUID" type="s" access="read"></property>
  <property name="Service" type="o" access="read"></property>
  <property name="Value" type="ay" access="read"></property>
  <property name="Notifying" type="b" access="read"></property>
  <property name="Flags" type="as" access="read"></property>
  <property name="Descriptors" type="ao" access="read"></property>
</interface>
`

//Device1IntrospectDataString interface definition
const Device1IntrospectDataString = `
<interface name="org.bluez.Device1">
  <method name="Disconnect"></method>
  <method name="Connect"></method>
  <method name="ConnectProfile">
    <arg name="UUID" type="s" direction="in"/>
  </method>
  <method name="DisconnectProfile">
    <arg name="UUID" type="s" direction="in"/>
  </method>
  <method name="Pair"></method>
  <method name="CancelPairing"></method>
  <property name="Address" type="s" access="read"></property>
  <property name="Name" type="s" access="read"></property>
  <property name="Alias" type="s" access="readwrite"></property>
  <property name="Class" type="u" access="read"></property>
  <property name="Appearance" type="q" access="read"></property>
  <property name="Icon" type="s" access="read"></property>
  <property name="Paired" type="b" access="read"></property>
  <property name="Trusted" type="b" access="readwrite"></property>
  <property name="Blocked" type="b" access="readwrite"></property>
  <property name="LegacyPairing" type="b" access="read"></property>
  <property name="RSSI" type="n" access="read"></property>
  <property name="Connected" type="b" access="read"></property>
  <property name="UUIDs" type="as" access="read"></property>
  <property name="Modalias" type="s" access="read"></property>
  <property name="Adapter" type="o" access="read"></property>
  <property name="ManufacturerData" type="a{qv}" access="read"></property>
  <property name="ServiceData" type="a{sv}" access="read"></property>
  <property name="TxPower" type="n" access="read"></property>
  <property name="GattServices" type="ao" access="read"></property>
</interface>
`

// GattService1IntrospectData interface definition
var GattService1IntrospectData = introspect.Interface{
	Name: "org.bluez.GattService1",
	Properties: []introspect.Property{
		{
			Name:   "UUID",
			Access: "read",
			Type:   "s",
		},
		{
			Name:   "Device",
			Access: "read",
			Type:   "o",
		},
		{
			Name:   "Primary",
			Access: "read",
			Type:   "b",
		},
		{
			Name:   "Characteristics",
			Access: "read",
			Type:   "ao",
		},
	},
}
