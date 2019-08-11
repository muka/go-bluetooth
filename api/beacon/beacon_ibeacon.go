package beacon

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
)

type BeaconIBeacon struct {
	Type          string
	ProximityUUID string
	Major         uint16
	Minor         uint16
	MeasuredPower uint16
}

// From Apple specifications
// Byte(s) 	Name 						Value 		Notes
// 0 				Flags[0] 				0x02 			See Bluetooth 4.0 Core Specification , Volume 3, Appendix C, 18.1.
// 1 				Flags[1] 				0x01 			See Bluetooth 4.0 Core Specification , Volume 3, Appendix C, 18.1.
// 2 				Flags[2] 				0x06 			See Bluetooth 4.0 Core Specification , Volume 3, Appendix C, 18.1.
// 3 				Length 					0x1A 			See Bluetooth 4.0 Core Specification
// 4 				Type 						0xFF 			See Bluetooth 4.0 Core Specification
// 5 				Company ID[0] 	0x4C 			Must not be used for any purposes not specified by Apple.
// 6 				Company ID[1] 	0x00 			Must not be used for any purposes not specified by Apple.
// ---- Bluez data starts here ----
// 7 				Beacon Type[0] 	0x02 			Must be set to 0x02 for all Proximity Beacons
// 8 				Beacon Type[1] 	0x15 			Must be set to 0x15 for all Proximity Beacons
// 9-24 		Proximity UUID 	0xnn..nn 	See CLBeaconRegion class in iOS Developer Library. Must not be set to all 0s.
// 25-26 		Major 					0xnnnn 		See CLBeaconRegion class in iOS Developer Library. 0x0000 = unset.
// 27-28 		Minor 					0xnnnn 		See CLBeaconRegion class in iOS Developer Library. 0x0000 = unset.
// 29 			Measured Power 	0xnn 			See Measured Power (page 7)
func (b *Beacon) ParseIBeacon(frames []uint8) BeaconIBeacon {

	info := BeaconIBeacon{}

	if frames[7-7] == 0x02 && frames[8-7] == 0x15 {
		info.Type = "proximity"
	}

	uuid := strings.ToUpper(hex.EncodeToString(frames[9-7 : 25-7]))
	info.ProximityUUID = strings.ToUpper(uuid)

	info.Major = binary.BigEndian.Uint16(frames[25-7 : 27-7])
	info.Minor = binary.BigEndian.Uint16(frames[27-7 : 29-7])

	if len(frames) < 23 {
		frames = append(frames, 0xb3)
	}

	info.MeasuredPower = uint16(frames[29-7])

	return info
}
