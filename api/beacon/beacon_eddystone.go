package beacon

import (
	"strings"

	log "github.com/sirupsen/logrus"
	eddystone "github.com/suapapa/go_eddystone"
)

const eddystoneSrvcUid = "FEAA"

type BeaconEddystone struct {
	Frame             eddystone.Header
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
	frameHeader := eddystone.Header(frames[0])

	switch frameHeader {
	case eddystone.UID:
		info.Frame = frameHeader
		parseEddystoneUID(&info, frames)
	case eddystone.TLM:
		info.Frame = frameHeader
		parseEddystoneTLM(&info, frames)
	case eddystone.URL:
		info.Frame = frameHeader
		err := parseEddystoneURL(&info, frames)
		if err != nil {
			log.Warn(err)
		}
	case eddystone.EID:
		// TODO

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
	ns, instance, tx := eddystone.ParseUIDFrame(frames)

	info.CalibratedTxPower = tx
	info.UID = strings.ToUpper(ns)
	info.InstanceUID = strings.ToUpper(instance)
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
	batt, temp, advCnt, secCnt := eddystone.ParseTLMFrame(frames)

	info.TLMVersion = int(frames[1] & 0xff)
	info.TLMBatteryVoltage = batt
	info.TLMTemperature = temp
	info.TLMAdvertisingPDU = advCnt
	info.TLMLastRebootedTime = secCnt
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
	url, tx, err := eddystone.ParseURLFrame(frames)
	if err != nil {
		return err
	}

	info.CalibratedTxPower = tx
	info.URL = url

	return nil
}
