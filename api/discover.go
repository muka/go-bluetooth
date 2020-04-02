package api

import (
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	log "github.com/sirupsen/logrus"
)

// Discover start device discovery
func Discover(
	a *adapter.Adapter1, filter *adapter.DiscoveryFilter,
) (
	chan *adapter.DeviceDiscovered, func(), error,
) {

	err := a.SetPairable(false)
	if err != nil {
		return nil, nil, err
	}
	err = a.SetDiscoverable(false)
	if err != nil {
		return nil, nil, err
	}
	err = a.SetPowered(true)
	if err != nil {
		return nil, nil, err
	}

	filterMap := make(map[string]interface{})
	if filter != nil {
		filterMap = filter.ToMap()
	}
	err = a.SetDiscoveryFilter(filterMap)
	if err != nil {
		return nil, nil, err
	}

	err = a.StartDiscovery()
	if err != nil {
		return nil, nil, err
	}

	ch, discoveryCancel, err := a.OnDeviceDiscovered()
	if err != nil {
		return nil, nil, err
	}

	cancel := func() {
		err := a.StopDiscovery()
		if err != nil {
			log.Warnf("Error stopping discovery: %s", err)
		}
		discoveryCancel()
	}

	return ch, cancel, nil
}
