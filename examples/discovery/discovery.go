//shows how to watch for new devices and list them
package discovery_example

import (
	"context"
	"os"
	"os/signal"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/api/beacon"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
	eddystone "github.com/suapapa/go_eddystone"
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

			go func(ev *adapter.DeviceDiscovered) {
				err = handleBeacon(dev)
				if err != nil {
					log.Errorf("%s: %s", ev.Path, err)
				}
			}(ev)
		}

	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill) // get notified of all OS signals

	sig := <-ch
	log.Infof("Received signal [%v]; shutting down...\n", sig)

	return nil
}

func handleBeacon(dev *device.Device1) error {

	b, err := beacon.NewBeacon(dev)
	if err != nil {
		return err
	}

	beaconUpdated, err := b.WatchDeviceChanges(context.Background())
	if err != nil {
		return err
	}

	isBeacon := <-beaconUpdated
	if !isBeacon {
		return nil
	}

	name := b.Device.Properties.Alias
	if name == "" {
		name = b.Device.Properties.Name
	}

	log.Debugf("Found beacon %s %s", b.Type, name)

	if b.IsEddystone() {
		ed := b.GetEddystone()
		switch ed.Frame {
		case eddystone.UID:
			log.Debugf(
				"Eddystone UID %s instance %s (%ddbi)",
				ed.UID,
				ed.InstanceUID,
				ed.CalibratedTxPower,
			)
			break
		case eddystone.TLM:
			log.Debugf(
				"Eddystone TLM temp:%.0f batt:%d last reboot:%d advertising pdu:%d (%ddbi)",
				ed.TLMTemperature,
				ed.TLMBatteryVoltage,
				ed.TLMLastRebootedTime,
				ed.TLMAdvertisingPDU,
				ed.CalibratedTxPower,
			)
			break
		case eddystone.URL:
			log.Debugf(
				"Eddystone URL %s (%ddbi)",
				ed.URL,
				ed.CalibratedTxPower,
			)
			break
		}

	}
	if b.IsIBeacon() {
		ibeacon := b.GetIBeacon()
		log.Debugf(
			"IBeacon %s (%ddbi) (major=%d minor=%d)",
			ibeacon.ProximityUUID,
			ibeacon.MeasuredPower,
			ibeacon.Major,
			ibeacon.Minor,
		)
	}

	return nil
}
