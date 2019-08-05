package api

import (
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	log "github.com/sirupsen/logrus"
)

// Discover start device discovery
func Discover(adapterID string, filter *adapter.DiscoveryFilter) (chan *adapter.DeviceDiscovered, func(), error) {

	a, err := adapter.NewAdapter1FromAdapterID(adapterID)
	if err != nil {
		return nil, nil, err
	}

	if filter != nil {
		err = a.SetDiscoveryFilter(filter.ToMap())
		if err != nil {
			return nil, nil, err
		}
	}

	err = a.StartDiscovery()
	if err != nil {
		return nil, nil, err
	}

	ch, discoveryCancel, err := a.DeviceDiscovered()

	cancel := func() {
		err := a.StopDiscovery()
		if err != nil {
			log.Warnf("Error stopping discovery: %s", err)
		}
		discoveryCancel()
	}

	return ch, cancel, nil
}
