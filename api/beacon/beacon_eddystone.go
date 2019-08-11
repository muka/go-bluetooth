package beacon

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
)

const eddystoneSrvcUid = "FEAA"

const (
	frameTypeUID byte = 0x00
	frameTypeURL      = 0x10
	frameTypeTLM      = 0x20
)

type EddystoneFrame string

const (
	EddystoneFrameUID EddystoneFrame = "uid"
	EddystoneFrameURL                = "url"
	EddystoneFrameTLM                = "tlm"
)

type BeaconEddystone struct {
	Frame             EddystoneFrame
	CalibratedTxPower int
	// eddystone-uid
	UID         string
	InstanceUID string

	URL string

	// eddystone-tlm plain
	TLMVersion          int
	TLMBatteryVoltage   uint16
	TLMTemperature      float32
	TLMAdvertisingPDU   uint32
	TLMLastRebootedTime uint32
}

func (b *Beacon) ParseEddystone(frames []byte) BeaconEddystone {

	info := BeaconEddystone{}

	switch frames[0] {
	case frameTypeUID:
		{
			info.Frame = EddystoneFrameUID
			parseEddystoneUID(&info, frames)
		}
		break
	case frameTypeTLM:
		{
			info.Frame = EddystoneFrameTLM
			parseEddystoneTLM(&info, frames)
		}
		break
	case frameTypeURL:
		{
			info.Frame = EddystoneFrameURL
			parseEddystoneURL(&info, frames)
		}
		break

	}

	return info
}

// eddystone-uid
// https://github.com/google/eddystone/tree/master/eddystone-uid
// Byte offset  Field	Description
// 0            Frame Type	Value = 0x00
// 1	          Ranging Data	Calibrated Tx power at 0 m
// 2	          NID[0]	10-byte Namespace
// 3	          NID[1]
// 4	          NID[2]
// 5	          NID[3]
// 6	          NID[4]
// 7	          NID[5]
// 8	          NID[6]
// 9	          NID[7]
// 10	          NID[8]
// 11	          NID[9]
// 12	          BID[0]	6-byte Instance
// 13	          BID[1]
// 14	          BID[2]
// 15	          BID[3]
// 16	          BID[4]
// 17	          BID[5]
// 18	          RFU	Reserved for future use, must be0x00
// 19	          RFU	Reserved for future use, must be0x00
func parseEddystoneUID(info *BeaconEddystone, frames []byte) {

	// 10 bytes length
	uid := hex.EncodeToString(frames[2:12])
	uid = strings.ToUpper(uid)

	// 6 bytes length
	iuid := hex.EncodeToString(frames[12:18])
	iuid = strings.ToUpper(iuid)

	// log.Debugf("%s - %s", uid, iuid)

	info.CalibratedTxPower = int(frames[1] & 0xff)
	info.UID = uid
	info.InstanceUID = iuid

}

// eddystone-tlm (plain)
// https://github.com/google/eddystone/blob/master/eddystone-tlm/tlm-plain.md
// Byte offset   Field	Description
// 0	           Frame Type	Value = 0x20
// 1	           Version	TLM version, value = 0x00
// 2	           VBATT[0]	Battery voltage, 1 mV/bit
// 3	           VBATT[1]
// 4	           TEMP[0]	Beacon temperature
// 5	           TEMP[1]
// 6	           ADV_CNT[0]	Advertising PDU count
// 7	           ADV_CNT[1]
// 8	           ADV_CNT[2]
// 9	           ADV_CNT[3]
// 10	           SEC_CNT[0]	Time since power-on or reboot
// 11	           SEC_CNT[1]
// 12	           SEC_CNT[2]
// 13	           SEC_CNT[3]
func parseEddystoneTLM(info *BeaconEddystone, frames []byte) {

	info.TLMVersion = int(frames[1] & 0xff)
	info.TLMBatteryVoltage = bytesToUint16(frames[2:4])
	info.TLMTemperature = fixTofloat32(bytesToUint16(frames[4:6]))
	info.TLMAdvertisingPDU = binary.BigEndian.Uint32(frames[6:10])
	info.TLMLastRebootedTime = binary.BigEndian.Uint32(frames[10:14])

}

// Byte offset	Field	Description
// 0	          Frame Type	Value = 0x10
// 1	          TX Power	Calibrated Tx power at 0 m
// 2	          URL Scheme	Encoded Scheme Prefix
// 3+	          Encoded URL	Length 1-17
//
// URL Scheme Prefix
// Decimal	 Hex   Expansion
// 0	       0x00	 http://www.
// 1	       0x01	 https://www.
// 2	       0x02	 http://
// 3	       0x03	 https://
//
// Eddystone-URL HTTP URL encoding
// Decimal	 Hex           Expansion
// 0	       0x00	         .com/
// 1	       0x01	         .org/
// 2	       0x02	         .edu/
// 3	       0x03	         .net/
// 4	       0x04	         .info/
// 5	       0x05	         .biz/
// 6	       0x06	         .gov/
// 7	       0x07	         .com
// 8	       0x08	         .org
// 9	       0x09	         .edu
// 10	       0x0a	         .net
// 11	       0x0b	         .info
// 12	       0x0c	         .biz
// 13	       0x0d	         .gov
// 14..32	   0x0e..0x20    Reserved for Future Use
// 127..255	 0x7F..0xFF    Reserved for Future Use
func parseEddystoneURL(info *BeaconEddystone, frames []byte) error {

	txPower := byteToInt(frames[1])
	info.CalibratedTxPower = txPower

	url, err := decodeURL(frames[2], frames[3:])
	if err != nil {
		return err
	}

	info.URL = url

	return nil
}
