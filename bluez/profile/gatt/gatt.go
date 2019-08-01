package gatt

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
