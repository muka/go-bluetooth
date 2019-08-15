package adapter

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

const (
	// DeviceRemoved a device has been removed from local cache
	DeviceRemoved uint8 = 0
	// DeviceAdded new device found, eg. via discovery
	DeviceAdded = iota
)

// DeviceDiscovered event emitted when a device is added or removed from Object Manager
type DeviceDiscovered struct {
	Path dbus.ObjectPath
	Type uint8
}

// OnDeviceDiscovered monitor for new devices and send updates via channel. Use cancel to close the monitoring process
func (a *Adapter1) OnDeviceDiscovered() (chan *DeviceDiscovered, func(), error) {

	signal, omSignalCancel, err := a.GetObjectManagerSignal()
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan *DeviceDiscovered)
	go (func() {
		for v := range signal {

			if v == nil {
				return
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

			path := v.Body[0].(dbus.ObjectPath)

			if op == DeviceRemoved {
				ifaces := v.Body[1].([]string)
				for _, iface := range ifaces {
					if iface == device.Device1Interface {
						log.Tracef("Removed device %s", path)
						ch <- &DeviceDiscovered{path, op}
					}
				}
				continue
			}

			ifaces := v.Body[1].(map[string]map[string]dbus.Variant)
			if p, ok := ifaces[device.Device1Interface]; ok {
				if p == nil {
					continue
				}
				log.Tracef("Added device %s", path)
				ch <- &DeviceDiscovered{path, op}
			}

		}
	})()

	cancel := func() {
		omSignalCancel()
		close(ch)
	}

	return ch, cancel, nil
}
