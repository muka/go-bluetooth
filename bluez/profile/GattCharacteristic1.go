package profile

import (
	"github.com/muka/go-bluetooth/src/gen/profile/gatt"
)

// Defines how the characteristic value can be used. See
// Core spec "Table 3.5: Characteristic Properties bit
// field", and "Table 3.8: Characteristic Extended
// Properties bit field"
const (
	FlagCharacteristicBroadcast                 = "broadcast"
	FlagCharacteristicRead                      = "read"
	FlagCharacteristicWriteWithoutResponse      = "write-without-response"
	FlagCharacteristicWrite                     = "write"
	FlagCharacteristicNotify                    = "notify"
	FlagCharacteristicIndicate                  = "indicate"
	FlagCharacteristicAuthenticatedSignedWrites = "authenticated-signed-writes"
	FlagCharacteristicReliableWrite             = "reliable-write"
	FlagCharacteristicWritableAuxiliaries       = "writable-auxiliaries"
	FlagCharacteristicEncryptRead               = "encrypt-read"
	FlagCharacteristicEncryptWrite              = "encrypt-write"
	FlagCharacteristicEncryptAuthenticatedRead  = "encrypt-authenticated-read"
	FlagCharacteristicEncryptAuthenticatedWrite = "encrypt-authenticated-write"
	FlagCharacteristicSecureRead                = "secure-read"
	FlagCharacteristicSecureWrite               = "secure-write"
)

// NewGattCharacteristic1 create a new GattCharacteristic1 client
func NewGattCharacteristic1(path string) (*gatt.GattCharacteristic1, error) {
	return gatt.NewGattCharacteristic1(path)
}
