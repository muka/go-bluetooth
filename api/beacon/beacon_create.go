package beacon

import (
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/muka/go-bluetooth/bluez/profile/advertising"
	eddystone "github.com/suapapa/go_eddystone"
)

func initBeacon() (*Beacon, error) {
	b := new(Beacon)
	b.props = new(advertising.LEAdvertisement1Properties)
	return b, nil
}

// CreateIBeacon Create a beacon in the IBeacon format
func CreateIBeacon(uuid string, major uint16, minor uint16, measuredPower uint16) (*Beacon, error) {

	frames := []byte{
		0x02, 0x15,
	}

	// uuid 2-17
	uuidBytes, err := hex.DecodeString(strings.Replace(uuid, "-", "", -1))
	if err != nil {
		return nil, err
	}
	frames = append(frames, uuidBytes...)

	// major 18,19
	mayorb := make([]byte, 2)
	binary.BigEndian.PutUint16(mayorb, major)
	frames = append(frames, mayorb...)

	// minor 20,21
	minorb := make([]byte, 2)
	binary.BigEndian.PutUint16(minorb, minor)
	frames = append(frames, minorb...)

	// pwr 22
	mpwr := make([]byte, 2)
	binary.BigEndian.PutUint16(mpwr, measuredPower)
	frames = append(frames, mpwr[1])

	b, err := initBeacon()
	if err != nil {
		return nil, err
	}

	b.Type = BeaconTypeIBeacon
	b.iBeacon = BeaconIBeacon{
		ProximityUUID: uuid,
		Major:         major,
		Minor:         minor,
		MeasuredPower: measuredPower,
		Type:          "proximity",
	}

	b.props.AddManifacturerData(appleBit, frames)

	return b, nil
}

func appendEddystoneService(UUIDs []string) []string {
	found := false
	for _, uuid := range UUIDs {
		if uuid == eddystoneSrvcUid {
			found = true
		}
	}
	if !found {
		return append(UUIDs, eddystoneSrvcUid)
	}
	return UUIDs
}

// CreateEddystoneURL create an eddystone beacon frame with url
func CreateEddystoneURL(url string, txPower int) (*Beacon, error) {

	frames, err := eddystone.MakeURLFrame(url, txPower)
	if err != nil {
		return nil, err
	}

	b, err := initBeacon()
	if err != nil {
		return nil, err
	}

	b.props.AddServiceUUID(eddystoneSrvcUid)
	b.props.AddServiceData(eddystoneSrvcUid, []byte(frames))
	b.Type = BeaconTypeEddystone
	b.eddystone = BeaconEddystone{
		URL:               url,
		CalibratedTxPower: txPower,
	}

	return b, nil
}

// CreateEddystoneTLM create an eddystone beacon frame with tlm
func CreateEddystoneTLM(batt uint16, temp float32, advCnt, secCnt uint32) (*Beacon, error) {

	frames, err := eddystone.MakeTLMFrame(batt, temp, advCnt, secCnt)
	if err != nil {
		return nil, err
	}

	b, err := initBeacon()
	if err != nil {
		return nil, err
	}

	b.props.AddServiceUUID(eddystoneSrvcUid)
	b.props.AddServiceData(eddystoneSrvcUid, []byte(frames))

	b.Type = BeaconTypeEddystone
	b.eddystone = BeaconEddystone{
		TLMVersion:          0,
		TLMTemperature:      temp,
		TLMAdvertisingPDU:   advCnt,
		TLMBatteryVoltage:   batt,
		TLMLastRebootedTime: secCnt,
	}

	return b, nil
}

// CreateEddystoneUID create an eddystone beacon frame with uid
func CreateEddystoneUID(namespace, instance string, txPwr int) (*Beacon, error) {

	frames, err := eddystone.MakeUIDFrame(namespace, instance, txPwr)
	if err != nil {
		return nil, err
	}

	b, err := initBeacon()
	if err != nil {
		return nil, err
	}

	b.props.AddServiceUUID(eddystoneSrvcUid)
	b.props.AddServiceData(eddystoneSrvcUid, []byte(frames))

	b.Type = BeaconTypeEddystone
	b.eddystone = BeaconEddystone{
		UID:               namespace,
		InstanceUID:       instance,
		CalibratedTxPower: txPwr,
	}

	return b, nil
}
