package service

import (
	"github.com/muka/go-bluetooth/bluez/profile/gatt"
)

// Create a new GattService1Properties
func NewGattService1Properties(uuid string) *gatt.GattService1Properties {
	return &gatt.GattService1Properties{
		IsService: true,
		Primary:   true,
		UUID:      uuid,
	}
}

// Create a new GattCharacteristic1Properties
func NewGattCharacteristic1Properties(uuid string) *gatt.GattCharacteristic1Properties {
	return &gatt.GattCharacteristic1Properties{
		UUID:  uuid,
		Flags: []string{},
	}
}

// Create a new GattDescriptor1Properties
func NewGattDescriptor1Properties(uuid string) *gatt.GattDescriptor1Properties {
	return &gatt.GattDescriptor1Properties{
		UUID: uuid,
		Flags: []string{
			gatt.FlagDescriptorRead,
			gatt.FlagDescriptorWrite,
		},
	}
}
