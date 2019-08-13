package beacon

import (
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/advertising"
)

// Expose the beacon
func (b *Beacon) Expose(adapterID string, timeout uint16) (func(), error) {

	props := b.props
	props.Type = advertising.AdvertisementTypeBroadcast

	if b.Name != "" {
		props.LocalName = b.Name
	}
	// Duration is set to 2sec by default
	// Not sure if duration can be mapped to interval.
	// props.Duration = 1
	props.Timeout = timeout

	cancel, err := api.ExposeAdvertisement(adapterID, props, uint32(timeout))

	return cancel, err
}
