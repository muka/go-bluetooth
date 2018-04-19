package profile

import (
	"github.com/fatih/structs"
)

// LEAdvertisement1Properties exposed properties for LEAdvertisement1
type LEAdvertisement1Properties struct {
	Type             string
	ServiceUUIDs     []string
	ManufacturerData map[string]interface{}
	//SolicitUUIDs     []string
	//ServiceData      map[string]interface{}
	//Includes         []string
	LocalName        string
	//Appearance       uint16
	//Duration         uint16
	//Timeout          uint16
}

//ToMap serialize properties
func (d *LEAdvertisement1Properties) ToMap() (map[string]interface{}, error) {
	return structs.Map(d), nil
}
