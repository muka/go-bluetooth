package zoeetree

import (
	"fmt"
	"time"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/emitter"
	log "github.com/sirupsen/logrus"
)

func NewZoeeTree(address string) *ZoeeTree {
	return &ZoeeTree{
		Address: address,
	}
}

type ZoeeTree struct {
	Address string
	Device  *api.Device
}

// Connect Attempt to connect in a defined timeout in seconds (0 means continuous)
func (z *ZoeeTree) Connect(timeout int) error {

	err := api.StartCleanDiscovery()
	if err != nil {
		return fmt.Errorf("Discovery failed: %s", err)
	}

	devch := make(chan *api.Device, 1)

	err = api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		dev := discoveryEvent.Device
		props, err := dev.GetProperties()
		if err != nil {
			log.Warn(fmt.Errorf("Failed to load properties for %s: %s", dev.Path, err))
			return
		}
		if props.Address == z.Address {
			devch <- dev
		}
	}))
	if err != nil {
		return fmt.Errorf("Discovery listener failed: %s", err)
	}

	if timeout > 0 {
		time.AfterFunc(time.Duration(timeout)*time.Second, func() {
			devch <- nil
		})
	}

	dev := <-devch

	if dev == nil {
		return fmt.Errorf("Device not found in %dsec", timeout)
	}

	z.Device = dev

	err = z.Device.Connect()
	if err != nil {
		return fmt.Errorf("Failed to connect: %s", err)
	}

	if !z.Device.Properties.Paired {
		err = z.Device.Pair()
		if err != nil {
			return fmt.Errorf("Failed to pair: %s", err)
		}
	}

	return nil
}
