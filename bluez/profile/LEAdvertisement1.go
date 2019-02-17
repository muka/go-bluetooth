package profile

import (
	"github.com/fatih/structs"
)

// LEAdvertisement1Properties exposed properties for LEAdvertisement1
type LEAdvertisement1Properties struct {
	Type             string
	ServiceUUIDs     []string
	ManufacturerData map[string]interface{}
	SolicitUUIDs     []string
	ServiceData      map[string]interface{}
	Data             map[string]interface{}
	// Advertise as general discoverable. When present this
	// will override adapter Discoverable property.
	//
	// Note: This property shall not be set when Type is set
	// to broadcast.
	Discoverable        bool
	DiscoverableTimeout uint16
	// List of features to be included in the advertising
	// packet.
	//
	// Possible values: as found on
	// 		LEAdvertisingManager.SupportedIncludes
	Includes   []string
	LocalName  string
	Appearance uint16
	Duration   uint16
	Timeout    uint16
}

//ToMap serialize properties
func (d *LEAdvertisement1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(d), nil
}
