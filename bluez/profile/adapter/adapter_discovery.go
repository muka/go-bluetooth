package adapter

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

const (
	DeviceRemoved uint8 = 0
	DeviceAdded         = iota
)

type DeviceDiscovered struct {
	Path dbus.ObjectPath
	Type uint8
}

func (a *Adapter1) DeviceDiscovered() (chan *DeviceDiscovered, error) {

	signal, err := a.GetPropertiesSignal()
	if err != nil {
		return nil, err
	}

	ch := make(chan *DeviceDiscovered)
	go (func() {
		for v := range signal {

			if v == nil {
				ch <- nil
				return
			}

			path := v.Body[0].(dbus.ObjectPath)
			ifaces := v.Body[1].(map[string]map[string]dbus.Variant)

			if p, ok := ifaces[device.Device1Interface]; ok {

				if p == nil {
					continue
				}

				var op uint8
				if v.Name == bluez.InterfacesAdded {
					op = DeviceAdded
				} else {
					if v.Name == bluez.InterfacesRemoved {
						op = DeviceRemoved
					} else {
						continue
					}
				}

				ch <- &DeviceDiscovered{path, op}
			}

		}
	})()

	return ch, nil
}
