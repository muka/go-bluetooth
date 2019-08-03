package adapter

import "github.com/muka/go-bluetooth/util"

const (
	DiscoveryFilterTransportAuto  = "auto"
	DiscoveryFilterTransportBrEdr = "bredr"
	DiscoveryFilterTransportLE    = "le"
)

type DiscoveryFilter struct {

	// Filter by service UUIDs, empty means match
	// _any_ UUID.
	// When a remote device is found that advertises
	// any UUID from UUIDs, it will be reported if:
	// - Pathloss and RSSI are both empty.
	// - only Pathloss param is set, device advertise
	// TX pwer, and computed pathloss is less than
	// Pathloss param.
	// - only RSSI param is set, and received RSSI is
	// higher than RSSI param.
	UUIDs []string

	// RSSI threshold value.
	// PropertiesChanged signals will be emitted
	// for already existing Device objects, with
	// updated RSSI value. If one or more discovery
	// filters have been set, the RSSI delta-threshold,
	// that is imposed by StartDiscovery by default,
	// will not be applied.
	RSSI int16

	// Pathloss threshold value.
	// PropertiesChanged signals will be emitted
	// for already existing Device objects, with
	// updated Pathloss value.
	Pathloss uint16

	// string Transport (Default "auto")
	// Transport parameter determines the type of
	// scan.
	// Possible values:
	// "auto"	- interleaved scan
	// "bredr"	- BR/EDR inquiry
	// "le"	- LE scan only
	// If "le" or "bredr" Transport is requested,
	// and the controller doesn't support it,
	// org.bluez.Error.Failed error will be returned.
	// If "auto" transport is requested, scan will use
	// LE, BREDR, or both, depending on what's
	// currently enabled on the controller.
	Transport string

	// bool DuplicateData (Default: true)
	// Disables duplicate detection of advertisement
	// data.
	DuplicateData bool
}

func (a *DiscoveryFilter) uuidExists(uuid string) bool {
	for _, uiid1 := range a.UUIDs {
		if uiid1 == uuid {
			return true
		}
	}
	return false
}

func (a *DiscoveryFilter) AddUUIDs(uuids ...string) {
	for _, uuid := range uuids {
		if !a.uuidExists(uuid) {
			a.UUIDs = append(a.UUIDs, uuid)
		}
	}
}

// ToMap convert to a format compatible with adapter SetDiscoveryFilter
func (a *DiscoveryFilter) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	util.StructToMap(a, m)
	return m
}

func NewDiscoveryFilter() DiscoveryFilter {
	return DiscoveryFilter{
		// default true
		DuplicateData: true,
	}
}
