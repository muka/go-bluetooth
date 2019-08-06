package rfkill

import (
	"errors"
	"strconv"

	"github.com/muka/go-bluetooth/hw/linux/hciconfig"
)

var rfclass = [...]string{
	"bluetooth",
	"wifi",
}

var rfkillHandler = NewRFKill()

// GetHCIConfig return an HCIConfig struct
func GetHCIConfig(adapterID string) *hciconfig.HCIConfig {
	return hciconfig.NewHCIConfig(adapterID)
}

// GetAdapterStatus return the status of an adapter
func GetAdapterStatus(adapterID string) (*RFKillResult, error) {

	if !rfkillHandler.IsInstalled() {
		return nil, errors.New("rfkill is not available")
	}

	list, err := rfkillHandler.ListAll()
	if err != nil {
		return nil, err
	}

	for _, adapter := range list {
		if adapter.Description == adapterID {
			// dbgSwitch("Got adapter index %d desc: %s type: %s hard-block: %t soft-block: %t",
			// 	adapter.Index,
			// 	adapter.Description,
			// 	adapter.IdentifierType,
			// 	adapter.HardBlocked,
			// 	adapter.SoftBlocked,
			// )

			return &adapter, nil
		}
	}

	return nil, errors.New("Adapter not found")
}

// ToggleAdapter Swap Off/On a device
func ToggleAdapter(adapterID string) error {
	err := TurnOffAdapter(adapterID)
	if err != nil {
		return err
	}
	return TurnOnAdapter(adapterID)
}

// TurnOnAdapter Enable a rfkill managed device
func TurnOnAdapter(adapterID string) error {

	var identifier string
	if isRFClass(adapterID) {
		identifier = adapterID
	} else {
		adapter, err := GetAdapterStatus(adapterID)
		if err != nil {
			return err
		}
		identifier = strconv.Itoa(adapter.Index)
	}

	if rfkillHandler.IsSoftBlocked(adapterID) {
		err := rfkillHandler.SoftUnblock(identifier)
		if err != nil {
			return err
		}
	}
	if rfkillHandler.IsHardBlocked(adapterID) {
		return errors.New("Adapter is hard locked, check for a physical switch to enable it")
	}
	return nil
}

// TurnOffAdapter Enable a rfkill managed device
func TurnOffAdapter(adapterID string) error {

	var identifier string
	if isRFClass(adapterID) {
		identifier = adapterID
	} else {
		adapter, err := GetAdapterStatus(adapterID)
		if err != nil {
			return err
		}
		identifier = strconv.Itoa(adapter.Index)
	}

	if !rfkillHandler.IsSoftBlocked(adapterID) {
		err := rfkillHandler.SoftBlock(identifier)
		if err != nil {
			return err
		}
	}
	return nil
}

func isRFClass(id string) bool {
	for _, class := range rfclass {
		if class == id {
			return true
		}
	}
	return false
}

// TurnOnBluetooth turn on bluetooth support
func TurnOnBluetooth() error {
	return TurnOnAdapter("bluetooth")
}

// TurnOffBluetooth turn on bluetooth support
func TurnOffBluetooth() error {
	return TurnOffAdapter("bluetooth")
}

// ToggleBluetooth toggle off/on the bluetooth support
func ToggleBluetooth() error {
	err := TurnOffBluetooth()
	if err != nil {
		return err
	}
	return TurnOnBluetooth()
}
