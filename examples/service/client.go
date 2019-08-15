package service_example

import (
	"errors"
	"fmt"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

func client(adapterID, hwaddr string) (err error) {

	log.Infof("Discovering %s on %s", hwaddr, adapterID)

	a, err := adapter.NewAdapter1FromAdapterID(adapterID)
	if err != nil {
		return err
	}

	dev, err := discover(a, hwaddr)
	if err != nil {
		return err
	}

	if dev == nil {
		return errors.New("Device not found, is it advertising?")
	}

	err = connect(dev)
	if err != nil {
		return err
	}

	watchProps, err := dev.WatchProperties()
	if err != nil {
		return err
	}

	for propUpdate := range watchProps {
		log.Debugf("propUpdate %++v", propUpdate)
	}

	return nil
}

func discover(a *adapter.Adapter1, hwaddr string) (*device.Device1, error) {

	err := a.FlushDevices()
	if err != nil {
		return nil, err
	}

	discovery, cancel, err := api.Discover(a, nil)
	if err != nil {
		return nil, err
	}

	defer cancel()

	for ev := range discovery {

		dev, err1 := device.NewDevice1(ev.Path)
		if err != nil {
			return nil, err1
		}

		if dev == nil || dev.Properties == nil {
			continue
		}

		p := dev.Properties

		n := p.Alias
		if p.Name != "" {
			n = p.Name
		}
		log.Debugf("Discovered (%s) %s", n, p.Address)

		if p.Address != hwaddr {
			continue
		}

		log.Infof("Found device %s", p.Address)
		return dev, nil
	}

	return nil, nil
}

func connect(dev *device.Device1) error {

	props := dev.Properties
	log.Infof("Found device name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)

	if props.Connected {
		return nil
	}

	if !props.Paired {
		log.Trace("Pairing device")
		err := dev.Pair()
		if err != nil {
			return fmt.Errorf("Pair failed: %s", err)
		}
	}

	log.Trace("Connecting device")
	err := dev.Connect()
	if err != nil {
		return fmt.Errorf("Connect failed: %s", err)
	}

	log.Debug("Device connected")
	return nil
}
