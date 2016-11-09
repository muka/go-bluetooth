package api

import (
	"errors"
	"strconv"

	"github.com/muka/bluez-client/linux"
)

var rfkill = linux.NewRFKill()

// ToggleDevice Swap Off/On a device
func ToggleDevice(adapterID string) error {
	identifier, err := GetRFKillAdapterIndex(adapterID)
	if err != nil {
		return err
	}
	err = TurnOffDevice(strconv.Itoa(identifier))
	if err != nil {
		return err
	}
	return TurnOnDevice(adapterID)
}

// TurnOnDevice Enable a rfkill managed device
func TurnOnDevice(adapterID string) error {
	identifier, err := GetRFKillAdapterIndex(adapterID)
	if err != nil {
		return err
	}
	if !rfkill.IsInstalled() {
		return errors.New("rfkill is not available")
	}
	if rfkill.IsSoftBlocked(adapterID) {
		err := rfkill.SoftUnblock(strconv.Itoa(identifier))
		if err != nil {
			return err
		}
	}
	if rfkill.IsHardBlocked(adapterID) {
		return errors.New("Device is hard locked, check for a physical switch to enable it")
	}
	return nil
}

// TurnOffDevice Enable a rfkill managed device
func TurnOffDevice(adapterID string) error {
	identifier, err := GetRFKillAdapterIndex(adapterID)
	if err != nil {
		return err
	}
	if !rfkill.IsInstalled() {
		return errors.New("rfkill is not available")
	}
	if !rfkill.IsSoftBlocked(adapterID) {
		err := rfkill.SoftBlock(strconv.Itoa(identifier))
		if err != nil {
			return err
		}
	}
	return nil
}

// GetRFKillAdapterIndex Return the adapter index from its name
func GetRFKillAdapterIndex(adapterID string) (int, error) {
	list, err := rfkill.ListAll()
	if err != nil {
		return -1, err
	}
	for _, adapter := range list {
		if adapter.Description == adapterID {
			logger.Debugf("Matching adapter index %d desc: %s type: %s hard-block: %t soft-block: %t",
				adapter.Index,
				adapter.Description,
				adapter.IdentifierType,
				adapter.HardBlocked,
				adapter.SoftBlocked,
			)

			return adapter.Index, nil
		}
	}
	return -1, errors.New("Adapter name not found")
}
