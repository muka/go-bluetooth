package profile

import (
	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

// Descriptor specific flags
const (
	FlagDescriptorRead                      = "read"
	FlagDescriptorWrite                     = "write"
	FlagDescriptorEncryptRead               = "encrypt-read"
	FlagDescriptorEncryptWrite              = "encrypt-write"
	FlagDescriptorEncryptAuthenticatedRead  = "encrypt-authenticated-read"
	FlagDescriptorEncryptAuthenticatedWrite = "encrypt-authenticated-write"
	FlagDescriptorSecureRead                = "secure-read"
	FlagDescriptorSecureWrite               = "secure-write"
)

// NewGattDescriptor1 create a new GattDescriptor1 client
func NewGattDescriptor1(objectPath string) (*gatt.GattDescriptor1, error) {
	return gatt.NewGattDescriptor1(objectPath)
}

// GattDescriptor1Properties exposed properties for GattDescriptor1
type GattDescriptor1Properties struct {
	Value          []byte `dbus:"emit"`
	Characteristic dbus.ObjectPath
	UUID           string
	Flags          []string
}
