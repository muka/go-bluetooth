//shows how to watch for new devices and list them
package discovery_example

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

func Run(adapterID string, onlyBeacon bool) error {

	//clean up connection on exit
	defer api.Exit()

	a, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	log.Debug("Flush cached devices")
	err = a.FlushDevices()
	if err != nil {
		return err
	}

	log.Debug("Start discovery")
	discovery, cancel, err := api.Discover(a, nil)
	if err != nil {
		return err
	}
	defer cancel()

	go func() {

		for ev := range discovery {

			if ev.Type == adapter.DeviceRemoved {
				continue
			}

			dev, err := device.NewDevice1(ev.Path)
			if err != nil {
				log.Errorf("%s: %s", ev.Path, err)
				continue
			}

			if dev == nil {
				log.Errorf("%s: not found", ev.Path)
				continue
			}

			log.Infof("name=%s addr=%s rssi=%d", dev.Properties.Name, dev.Properties.Address, dev.Properties.RSSI)

			err = handleBeacon(dev)
			if err != nil {
				log.Errorf("%s: %s", ev.Path, err)
			}
		}

	}()

	select {}
}

func showDeviceInfo(path dbus.ObjectPath, onlyBeacon bool) error {

	return nil
}

func handleBeacon(dev *device.Device1) error {

	isBeacon, b, err := api.NewBeacon(dev)
	if err != nil {
		return err
	}

	if !isBeacon {
		return nil
	}

	log.Debugf("Found beacon %s %s", b.Type, b.Device.Properties.Name)

	return nil
}
