package mesh

import "github.com/godbus/dbus/v5"

//VendorItem array{(uint16, uint16)}
type VendorItem struct {
	Vendor  uint16
	ModelID uint16
}

//VendorItem array{(uint16, uint16, dict)}
type VendorOptionsItem struct {
	VendorItem
	Options map[string]interface{}
}

// ModelConfig
type ModelConfig struct {
	Bindings          []uint16
	PublicationPeriod uint32
	Vendor            uint16
	Subscriptions     []dbus.Variant
}

// Model
type Model struct {
	Identifier uint16
	Config     ModelConfig
}

// ConfigurationItem array{byte, array{(uint16, dict)}}
type ConfigurationItem struct {
	Index  byte
	Models []Model
}
