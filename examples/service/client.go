package service_example

import (
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

func client(adapterID, hwaddr string) (err error) {

	log.Infof("Discovering %s on %s", hwaddr, adapterID)

	a, err := adapter.NewAdapter1FromAdapterID(adapterID)
	if err != nil {
		return err
	}

	//Connect DBus System bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	// do not reuse agent0 from service
	agent.NextAgentPath()

	ag := agent.NewSimpleAgent()
	err = agent.ExposeAgent(conn, ag, agent.CapNoInputNoOutput, true)
	if err != nil {
		return fmt.Errorf("SimpleAgent: %s", err)
	}

	dev, err := discover(a, hwaddr)
	if err != nil {
		return err
	}
	if dev == nil {
		return errors.New("Device not found, is it advertising?")
	}

	watchProps, err := dev.WatchProperties()
	if err != nil {
		return err
	}
	go func() {
		for propUpdate := range watchProps {
			log.Debugf("> property updated %s", propUpdate.Name)
		}
	}()

	err = connect(dev, ag, adapterID)
	if err != nil {
		return err
	}

	retrieveServices(a, dev)

	select {}
	// return nil
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

		return dev, nil
	}

	return nil, nil
}

func connect(dev *device.Device1, ag *agent.SimpleAgent, adapterID string) error {

	props, err := dev.GetProperties()
	if err != nil {
		return fmt.Errorf("Failed to load props: %s", err)
	}

	log.Infof("Found device name=%s addr=%s rssi=%d", props.Name, props.Address, props.RSSI)

	if props.Connected {
		log.Trace("Device is connected")
		return nil
	}

	if !props.Paired || !props.Trusted {
		log.Trace("Pairing device")

		err := dev.Pair()
		if err != nil {
			return fmt.Errorf("Pair failed: %s", err)
		}

		log.Info("Pair succeed, connecting...")
		agent.SetTrusted(adapterID, dev.Path())
	}

	if !props.Connected {
		log.Trace("Connecting device")
		err = dev.Connect()
		if err != nil {
			if !strings.Contains(err.Error(), "Connection refused") {
				return fmt.Errorf("Connect failed: %s", err)
			}
		}
	}

	return nil
}

func retrieveServices(a *adapter.Adapter1, dev *device.Device1) error {

	log.Debug("Exposed services")
	list, err := dev.GetAllServicesAndUUID()
	if err != nil {
		return err
	}

	for _, servicePath := range list {
		log.Debugf("%s", servicePath)
	}

	return nil
}
