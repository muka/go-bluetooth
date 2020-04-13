package advertising

const (
	SecondaryChannel1M    = "1M"
	SecondaryChannel2M    = "2M"
	SecondaryChannelCoded = "Coded"
)

// "tx-power"
// "appearance"
// "local-name"
const (
	SupportedIncludesTxPower    = "tx-power"
	SupportedIncludesAppearance = "appearance"
	SupportedIncludesLocalName  = "local-name"
)

const (
	AdvertisementTypeBroadcast  = "broadcast"
	AdvertisementTypePeripheral = "peripheral"
)

func (a *LEAdvertisement1Properties) AddServiceUUID(uuids ...string) {
	if a.ServiceUUIDs == nil {
		a.ServiceUUIDs = make([]string, 0)
	}
	for _, uuid := range uuids {
		for _, uuid1 := range a.ServiceUUIDs {
			if uuid1 == uuid {
				continue
			}
		}
		a.ServiceUUIDs = append(a.ServiceUUIDs, uuid)
	}
}

func (a *LEAdvertisement1Properties) AddData(code byte, data []uint8) {
	if a.Data == nil {
		a.Data = make(map[byte]interface{})
	}
	a.Data[code] = data
}

func (a *LEAdvertisement1Properties) AddServiceData(code string, data []uint8) {
	if a.ServiceData == nil {
		a.ServiceData = make(map[string]interface{})
	}
	a.ServiceData[code] = data
}

func (a *LEAdvertisement1Properties) AddManifacturerData(code uint16, data []uint8) {
	if a.ManufacturerData == nil {
		a.ManufacturerData = make(map[uint16]interface{})
	}
	a.ManufacturerData[code] = data
}
