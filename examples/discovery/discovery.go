//shows how to watch for new devices and list them
package discovery_example

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

func Run(adapterID string) error {

	//clean up connection on exit
	defer api.Exit()

	a, err := adapter.GetAdapter(adapterID)
	if err != nil {
		return err
	}

	err = a.FlushDevices()
	if err != nil {
		return err
	}

	log.Debug("Start discovery")
	discovery, cancel, err := api.Discover(adapterID, nil)
	if err != nil {
		return err
	}
	defer cancel()

	go func() {

		for ev := range discovery {

			if ev.Type == adapter.DeviceRemoved {
				continue
			}

			err = showDeviceInfo(ev.Path)
			if err != nil {
				log.Errorf("Error: %s", err)
			}
		}

	}()

	select {}
}

func showDeviceInfo(path dbus.ObjectPath) error {

	dev, err := device.NewDevice1(path)
	if err != nil {
		return err
	}

	if dev == nil {
		return fmt.Errorf("Device not found %s", path)
	}
	props, err := dev.GetProperties()
	if err != nil {
		return fmt.Errorf("%s: Failed to get properties: %s", dev.Path(), err.Error())
	}

	log.Infof("name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)
	return nil
}
