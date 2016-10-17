package bluetooth

import (
	"github.com/godbus/dbus"
	"github.com/muka/bluez-client/bluez"
	"github.com/muka/device-manager/util"
)

var manager = bluez.NewObjectManager()

//GetDevices returns a list of bluetooth discovered Devices
func GetDevices() ([]Device, error) {

	objects, err := manager.GetManagedObjects()

	if err != nil {
		return nil, err
	}

	var devices = make([]Device, 0)
	for path, list := range objects {
		for iface, variant := range list {
			switch iface {
			case "org.bluez.Device1":
				{
					deviceProperties := new(bluez.Device1Properties)
					util.MapToStruct(deviceProperties, variant.Value().(map[string]dbus.Variant))
					dev := NewDevice(string(path), deviceProperties)
					devices = append(devices, *dev)
				}
			}
		}
	}

	return devices, nil
}
