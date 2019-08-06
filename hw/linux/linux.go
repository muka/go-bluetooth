package linux

import (
	"fmt"
	"strconv"

	"github.com/muka/go-bluetooth/hw/linux/btmgmt"
	"github.com/muka/go-bluetooth/hw/linux/hci"
	"github.com/muka/go-bluetooth/hw/linux/hciconfig"
	log "github.com/sirupsen/logrus"
)

type BackendType string

const (
	BackendBtmgmt    BackendType = "btmgmt"
	BackendHCI       BackendType = "hci"
	BackendHCIConfig BackendType = "hciconfig"
)

var Backend BackendType = BackendHCIConfig

type AdapterInfo struct {
	AdapterID string
	Address   string
	Type      string
	Enabled   bool
}

// GetAdapter return status information for a controller
func GetAdapter(adapterID string) (a AdapterInfo, err error) {

	list, err := GetAdapters()
	if err != nil {
		return a, err
	}

	for _, a := range list {
		if a.AdapterID == adapterID {
			return a, err
		}
	}

	return a, fmt.Errorf("Adapter %s not found", adapterID)
}

// GetAdapters return a list of status information of available controllers
func GetAdapters() ([]AdapterInfo, error) {

	list, err := hciconfig.GetAdapters()
	if err != nil {
		return nil, err
	}

	list1 := []AdapterInfo{}
	for _, info := range list {
		list1 = append(list1, AdapterInfo{
			AdapterID: info.AdapterID,
			Enabled:   info.Enabled,
			Type:      info.Type,
			Address:   info.Address,
		})
	}

	return list1, err
}

func Up(adapterID string) error {

	status, err := GetAdapter(adapterID)
	if err != nil {
		return err
	}

	if status.Enabled {
		return nil
	}

	if Backend == BackendHCIConfig {
		_, err := hciconfig.NewHCIConfig(adapterID).Up()
		return err
	}

	if Backend == BackendBtmgmt {
		return btmgmt.NewBtMgmt(adapterID).SetPowered(true)
	}

	if Backend == BackendHCI {

		id, err := strconv.Atoi(adapterID[3:])
		if err != nil {
			return err
		}
		return hci.Up(id)
	}

	return fmt.Errorf("Unsupported backend type: %s", Backend)
}

func Down(adapterID string) error {

	status, err := GetAdapter(adapterID)
	if err != nil {
		return err
	}

	if !status.Enabled {
		return nil
	}

	if Backend == BackendHCIConfig {
		_, err := hciconfig.NewHCIConfig(adapterID).Down()
		return err
	}

	if Backend == BackendBtmgmt {
		return btmgmt.NewBtMgmt(adapterID).SetPowered(false)
	}

	if Backend == BackendHCI {
		id, err := strconv.Atoi(adapterID[3:])
		if err != nil {
			return err
		}
		return hci.Down(id)
	}

	return fmt.Errorf("Unsupported backend type: %s", Backend)
}

func Reset(adapterID string) error {
	err := Down(adapterID)
	if err != nil {
		log.Warnf("Down failed: %s", err)
	}
	return Up(adapterID)
}
